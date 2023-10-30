package handlers

import (
	"go-metrics/internal/server/storage"
)

type Ping interface {
	Ping() error
}

type MuxHandlers struct {
	s storage.Storage
	p Ping
}

func NewMuxHandlers(s storage.Storage, p Ping) *MuxHandlers {
	return &MuxHandlers{
		s: s,
		p: p,
	}
}
