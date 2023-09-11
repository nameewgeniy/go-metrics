package memory

import (
	"sync"
)

type MemStorage struct {
	gauge   sync.Map
	counter sync.Map
}

func (i *MemStorage) AddGauge(name string, value float64) {
	i.gauge.Store(name, value)
}

func (i *MemStorage) AddCount(name string, value int64) {
	if oldValue, ok := i.counter.Load(name); ok {
		i.counter.Store(name, oldValue.(int64)+value)
	}
}
