package memory

import (
	"go-metrics/internal/server/storage"
)

func (m *Memory) AddCounter(counter storage.MetricsItemCounter) error {
	if oldValue, ok := m.Counter.Load(counter.Name); ok {
		counter.Value += oldValue.(int64)
	}
	m.Counter.Store(counter.Name, counter.Value)

	return nil
}

func (m *Memory) FindCounterItem(name string) (storage.MetricsItemCounter, error) {
	res := storage.MetricsItemCounter{}
	if val, ok := m.Counter.Load(name); ok {
		res.Value = val.(int64)
		return res, nil
	}
	return storage.MetricsItemCounter{}, storage.ErrItemNotFound
}

func (m *Memory) FindCounterAll() ([]storage.MetricsItemCounter, error) {
	var res []storage.MetricsItemCounter

	m.Counter.Range(func(name, value interface{}) bool {
		res = append(res, storage.MetricsItemCounter{
			Name:  name.(string),
			Value: value.(int64),
		})
		return true
	})

	return res, nil
}
