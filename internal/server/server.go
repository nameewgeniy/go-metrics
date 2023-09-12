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
	UpdateMetricsHandle(http.ResponseWriter, *http.Request)
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
	r.HandleFunc("/update/{type}/{name}/{value}", s.h.UpdateMetricsHandle)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.cnf.Addr(),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv.ListenAndServe()
}
