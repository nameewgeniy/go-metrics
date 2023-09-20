package memory

import (
	"sync"
)

type Memory struct {
	Gauge   sync.Map
	Counter sync.Map
}

func NewMemory() *Memory {
	return &Memory{}
}
