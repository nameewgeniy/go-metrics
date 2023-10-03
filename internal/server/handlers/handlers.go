package handlers

import (
	"go-metrics/internal/server/storage"
)

type MuxHandlers struct {
	s storage.Storage
}

func NewMuxHandlers(s storage.Storage) *MuxHandlers {
	return &MuxHandlers{
		s: s,
	}
}
