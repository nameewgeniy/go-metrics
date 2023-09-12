package server

import (
	"net/http"
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

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, s.h.UpdateMetricsHandle)

	return http.ListenAndServe(s.cnf.Addr(), mux)
}
