package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/middleware"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/memory"
	"net/http"
	"os"
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

func (s Server) Listen(_ context.Context, errorCh chan<- error) {

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

	if err := srv.ListenAndServe(); err != nil {
		errorCh <- err
	}
}

func (s Server) Workers(ctx context.Context, errorCh chan<- error, sig chan os.Signal) {

	defer func() {
		if r := recover(); r != nil {
			errorCh <- fmt.Errorf("snapshot panic: %v", r)
		}
	}()

	// Запускаем интервальный сброс данных в файл
	snapshotTicker := time.NewTicker(s.cnf.StoreInterval())
	defer snapshotTicker.Stop()

	for {
		select {
		case <-sig:
			_ = s.sn.Snapshot()
			return
		case <-ctx.Done():
			_ = s.sn.Snapshot()
			return
		case <-snapshotTicker.C:
			if err := s.sn.Snapshot(); err != nil {
				errorCh <- err
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
