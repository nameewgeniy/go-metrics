package memory

import (
	"encoding/json"
	"os"
)

func (m *Memory) Snapshot() error {

	counters, err := m.FindCounterAll()
	if err != nil {
		return err
	}

	gauges, err := m.FindGaugeAll()
	if err != nil {
		return err
	}

	if len(gauges) == 0 && len(counters) == 0 {
		return nil
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
