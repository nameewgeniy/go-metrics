package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ServerConfig interface {
	Addr() string
}

type Handlers interface {
	UpdateGauageMetricsHandle(http.ResponseWriter, *http.Request)
	GetGauageMetricsHandle(http.ResponseWriter, *http.Request)

	UpdateCounterMetricsHandle(http.ResponseWriter, *http.Request)
	GetCounterMetricsHandle(http.ResponseWriter, *http.Request)

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
	r.HandleFunc("/", s.h.ViewMetricsHandle).Methods(http.MethodGet)

	r.HandleFunc("/update/gauage/{name}/{value}", s.h.UpdateGauageMetricsHandle).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/value/gauage/{name}", s.h.GetGauageMetricsHandle).Methods(http.MethodGet)

	r.HandleFunc("/update/counter/{name}/{value}", s.h.UpdateCounterMetricsHandle).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/value/counter/{name}", s.h.GetCounterMetricsHandle).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.cnf.Addr(),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv.ListenAndServe()
}
