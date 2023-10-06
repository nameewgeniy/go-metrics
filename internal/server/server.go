package server

import (
	"context"
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/middleware"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/memory"
	"net/http"
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

func (s Server) Listen(ctx context.Context) error {

	r := mux.NewRouter()
	r.Handle("/", middleware.RequestLogger(middleware.CompressMiddleware(s.h.ViewMetricsHandle))).Methods(http.MethodGet)

	r.Handle("/update/", middleware.RequestLogger(middleware.CompressMiddleware(s.h.UpdateMetricsJSONHandle))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/value/", middleware.RequestLogger(middleware.CompressMiddleware(s.h.GetMetricsJSONHandle))).Methods(http.MethodPost, http.MethodOptions)

	r.Handle("/update/{type}/{name}/{value}", middleware.RequestLogger(middleware.CompressMiddleware(s.h.UpdateMetricsHandle))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/value/{type}/{name}", middleware.RequestLogger(middleware.CompressMiddleware(s.h.GetMetricsHandle))).Methods(http.MethodGet)

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

func (s Server) Snapshot(ctx context.Context) error {

	// Запускаем интервальный сброс данных в файл
	snapshotTicker := time.NewTicker(s.cnf.StoreInterval())
	defer snapshotTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return s.sn.Snapshot()
		case <-snapshotTicker.C:
			if err := s.sn.Snapshot(); err != nil {
				return err
			}
		}
	}
}

// Restore Выгружаем данные в storage
func (s Server) Restore() error {
	if s.cnf.Restore() {
		return s.sn.Restore()
	}

	return nil
}
