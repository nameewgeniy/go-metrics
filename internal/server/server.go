package server

import (
	"github.com/gorilla/mux"
	"github.com/nameewgeniy/go-metrics/internal/logger/middleware"
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
	r.Handle("/", middleware.RequestLogger(s.h.ViewMetricsHandle)).Methods(http.MethodGet)

	r.Handle("/update/{type}/{name}/{value}", middleware.RequestLogger(s.h.UpdateMetricsHandle)).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/value/{type}/{name}", middleware.RequestLogger(s.h.GetMetricsHandle)).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.cnf.Addr(),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv.ListenAndServe()
}
