package memory

import (
	"go-metrics/internal/server/storage"
	"sync"
)

type MemoryStorageConfig interface {
	FileStoragePath() string
	IsRestore() bool
}

type SnapshotStorage interface {
	Restore() error
	Snapshot() error
}

type Memory struct {
	Gauge   sync.Map
	Counter sync.Map
	cfg     MemoryStorageConfig
}

type snapshotItems struct {
	Gauge   []storage.MetricsItemGauge   `json:"gauge"`
	Counter []storage.MetricsItemCounter `json:"counter"`
}

var mutex sync.Mutex

func NewMemoryStorage(cfg MemoryStorageConfig) *Memory {
	return &Memory{
		cfg: cfg,
	}
}
