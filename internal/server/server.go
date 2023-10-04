package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/middleware"
	"go-metrics/internal/server/storage"
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
}

func NewServer(c ServerConfig, h Handlers, s storage.Storage) *Server {
	return &Server{
		cnf: c,
		h:   h,
		s:   s,
	}
}

func (s Server) Listen(errorCh chan<- error) {

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

func (s Server) Workers(errorCh chan<- error, sig chan os.Signal) {

	defer func() {
		if r := recover(); r != nil {
			errorCh <- fmt.Errorf("snapshot panic: %v", r)
		}
	}()

	// Выгружаем данные в storage
	if s.cnf.Restore() {
		if err := s.s.Restore(); err != nil {
			errorCh <- err
		}
	}

	// Запускаем интервальный сброс данных в файл
	snapshotTicker := time.NewTicker(s.cnf.StoreInterval())
	defer snapshotTicker.Stop()

	for {
		select {
		case <-sig:
			_ = s.s.Snapshot()
		case <-snapshotTicker.C:
			if err := s.s.Snapshot(); err != nil {
				errorCh <- err
			}
		}
	}
}
