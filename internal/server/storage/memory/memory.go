package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"sync"
)

type Memory struct {
	gauge   sync.Map
	counter sync.Map
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Save(i storage.MetricsItem) error {
	for name, value := range i.Gauge {
		m.AddGauge(name, value)
	}

	for name, value := range i.Counter {
		m.AddCount(name, value)
	}

	return nil
}

func (m *Memory) AddGauge(name string, value float64) {
	m.gauge.Store(name, value)
}

func (m *Memory) AddCount(name string, value int64) {
	if oldValue, ok := m.counter.Load(name); ok {
		m.counter.Store(name, oldValue.(int64)+value)
	}
}
