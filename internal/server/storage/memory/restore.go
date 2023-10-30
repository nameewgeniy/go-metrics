package memory

import (
	"encoding/json"
	"io"
	"os"
)

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
