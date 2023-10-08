package handlers

import (
	"go-metrics/internal/server/storage"
)

type PingDB interface {
	Ping() error
}

type MuxHandlers struct {
	s storage.Storage
	p PingDB
}

func NewMuxHandlers(s storage.Storage, p PingDB) *MuxHandlers {
	return &MuxHandlers{
		s: s,
		p: p,
	}
}
