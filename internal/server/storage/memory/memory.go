package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
)

type Memory struct {
	s *MemStorage
}

func NewMemory() *Memory {
	return &Memory{
		s: &MemStorage{
			gauge:   map[string]float64{},
			counter: map[string]int64{},
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
