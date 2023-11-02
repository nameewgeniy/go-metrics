package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/middleware"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/memory"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig interface {
	Addr() string
	StoreInterval() time.Duration
	Restore() bool
}

type Handlers interface {
	UpdateMetricsHandle(http.ResponseWriter, *http.Request)
	GetMetricsHandle(http.ResponseWriter, *http.Request)
	ViewMetricsHandle(http.ResponseWriter, *http.Request)

	GetMetricsJSONHandle(http.ResponseWriter, *http.Request)
	UpdateMetricsJSONHandle(http.ResponseWriter, *http.Request)

	PingHandle(http.ResponseWriter, *http.Request)

	UpdateBatchMetricsHandle(http.ResponseWriter, *http.Request)
}

type Server struct {
	cnf ServerConfig
	h   Handlers
	s   storage.Storage
	sn  memory.SnapshotStorage
}

func NewServer(c ServerConfig, h Handlers, s storage.Storage, sn memory.SnapshotStorage) *Server {
	return &Server{
		cnf: c,
		h:   h,
		s:   s,
		sn:  sn,
	}
}

func (s Server) Run() error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errorCh := make(chan error)
	defer close(errorCh)

	// Storage up
	if err := s.s.Up(ctx); err != nil {
		return err
	}

	// Storage down
	defer func() {
		_ = s.s.Down(ctx)
	}()

	eg, errCtx := errgroup.WithContext(ctx)

	// Interval snapshot run
	if s.sn != nil {
		eg.Go(func() error {
			defer handlePanic(errorCh, cancel)
			return s.intervalSnapshot(errCtx)
		})
	}

	// Server listen
	eg.Go(func() error {
		defer handlePanic(errorCh, cancel)
		return s.listen(errCtx)
	})

	go func() {
		errorCh <- eg.Wait()
	}()

	return <-errorCh
}

func (s Server) listen(ctx context.Context) error {

	r := mux.NewRouter()
	r.Handle("/", middleware.RequestLogger(middleware.CompressMiddleware(s.h.ViewMetricsHandle))).Methods(http.MethodGet)

	r.Handle("/update/", middleware.RequestLogger(middleware.CompressMiddleware(middleware.CheckSignature(s.h.UpdateMetricsJSONHandle)))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/value/", middleware.RequestLogger(middleware.CompressMiddleware(s.h.GetMetricsJSONHandle))).Methods(http.MethodPost, http.MethodOptions)

	r.Handle("/update/{type}/{name}/{value}", middleware.RequestLogger(middleware.CompressMiddleware(s.h.UpdateMetricsHandle))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/value/{type}/{name}", middleware.RequestLogger(middleware.CompressMiddleware(s.h.GetMetricsHandle))).Methods(http.MethodGet)

	r.Handle("/ping", middleware.RequestLogger(middleware.CompressMiddleware(s.h.PingHandle))).Methods(http.MethodGet)

	r.Handle("/updates/", middleware.RequestLogger(middleware.CompressMiddleware(middleware.CheckSignature(s.h.UpdateBatchMetricsHandle)))).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.cnf.Addr(),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s Server) intervalSnapshot(ctx context.Context) error {

	// Запускаем интервальный сброс данных в файл
	snapshotTicker := time.NewTicker(s.cnf.StoreInterval())
	defer snapshotTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-snapshotTicker.C:
			if err := s.sn.Snapshot(); err != nil {
				return err
			}
		}
	}
}

func handlePanic(errorCh chan<- error, stop context.CancelFunc) {
	if r := recover(); r != nil {
		errorCh <- fmt.Errorf("panic: %v", r)
		stop()
	}
}
