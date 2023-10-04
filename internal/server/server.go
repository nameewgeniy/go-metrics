package server

import (
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/middleware"
	"net/http"
	"time"
)

type ServerConfig interface {
	Addr() string
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
}

func NewServer(c ServerConfig, h Handlers) *Server {
	return &Server{
		cnf: c,
		h:   h,
	}
}

func (s Server) Listen() error {

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

	return srv.ListenAndServe()
}
