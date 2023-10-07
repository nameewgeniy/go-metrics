package handlers

import (
	"go-metrics/internal/server/storage"
)

type PingDb interface {
	Ping() error
}

type MuxHandlers struct {
	s storage.Storage
	p PingDb
}

func NewMuxHandlers(s storage.Storage, p PingDb) *MuxHandlers {
	return &MuxHandlers{
		s: s,
		p: p,
	}
}
