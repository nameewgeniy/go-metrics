package memory

import (
	"go-metrics/internal/server/storage"
)

func (m *Memory) AddGauge(gauge storage.MetricsItemGauge) error {
	m.Gauge.Store(gauge.Name, gauge.Value)

	return nil
}

func (m *Memory) FindGaugeItem(name string) (storage.MetricsItemGauge, error) {
	res := storage.MetricsItemGauge{}
	if val, ok := m.Gauge.Load(name); ok {
		res.Value = val.(float64)
		return res, nil
	}
	return storage.MetricsItemGauge{}, storage.ErrItemNotFound
}

func (m *Memory) FindGaugeAll() ([]storage.MetricsItemGauge, error) {
	var res []storage.MetricsItemGauge

	m.Gauge.Range(func(name, value interface{}) bool {
		res = append(res, storage.MetricsItemGauge{
			Name:  name.(string),
			Value: value.(float64),
		})
		return true
	})

	return res, nil
}
