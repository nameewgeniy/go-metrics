package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"sync"
)

type Memory struct {
	s *MemStorage
}

func NewMemory() *Memory {
	return &Memory{
		s: &MemStorage{
			gauge:   sync.Map{},
			counter: sync.Map{},
		},
	}
}

func (m Memory) Save(i storage.MetricsItem) error {
	for name, value := range i.Gauge {
		m.s.AddGauge(name, value)
	}

	for name, value := range i.Counter {
		m.s.AddCount(name, value)
	}

	return nil
}
