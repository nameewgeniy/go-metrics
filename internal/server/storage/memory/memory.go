package memory

import (
	"encoding/json"
	"go-metrics/internal/server/storage"
	"io"
	"os"
	"sync"
)

type StorageConfig interface {
	FileStoragePath() string
}

type Memory struct {
	Gauge   sync.Map
	Counter sync.Map
	cfg     StorageConfig
}

type snapshotItems struct {
	Gauge   []storage.MetricsItemGauge   `json:"gauge"`
	Counter []storage.MetricsItemCounter `json:"counter"`
}

var mutex sync.Mutex

func NewMemory(cfg StorageConfig) *Memory {
	return &Memory{
		cfg: cfg,
	}
}

func (m *Memory) Snapshot() error {

	counters, err := m.FindCounterAll()
	if err != nil {
		return err
	}

	gauges, err := m.FindGaugeAll()
	if err != nil {
		return err
	}

	items := snapshotItems{
		Gauge:   gauges,
		Counter: counters,
	}

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return m.writeToFile(data)
}

func (m *Memory) writeToFile(data []byte) error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(m.cfg.FileStoragePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return err
	}

	return nil
}

func (m *Memory) Restore() error {

	data, err := m.readFromFile()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	items := snapshotItems{}
	if err = json.Unmarshal(data, &items); err != nil {
		return err
	}

	for _, v := range items.Counter {
		m.Counter.Store(v.Name, v.Value)
	}

	for _, v := range items.Gauge {
		m.Gauge.Store(v.Name, v.Value)
	}

	return nil
}

func (m *Memory) readFromFile() ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.OpenFile(m.cfg.FileStoragePath(), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
