package server

import (
	"net/http"
)

type Conf interface {
	Addr() string
}

type Handlers interface {
	UpdateMetricsHande(http.ResponseWriter, *http.Request)
}

type Server struct {
	cnf Conf
	h   Handlers
}

func NewServer(c Conf, h Handlers) *Server {
	return &Server{
		cnf: c,
		h:   h,
	}
}

func (s Server) Listen() error {

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, s.h.UpdateMetricsHande)

	return http.ListenAndServe(s.cnf.Addr(), mux)
}
