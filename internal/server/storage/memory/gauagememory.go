package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
)

func (m *Memory) AddGauage(gauage storage.MetricsItemGauage) error {
	m.Gauge.Store(gauage.Name, gauage.Value)

	return nil
}

func (m *Memory) FindGauageItem(name string) (storage.MetricsItemGauage, error) {
	res := storage.MetricsItemGauage{}
	if val, ok := m.Gauge.Load(name); ok {
		res.Value = val.(float64)
	}
	return res, nil
}

func (m *Memory) FindGauageAll() ([]storage.MetricsItemGauage, error) {
	var res []storage.MetricsItemGauage

	m.Gauge.Range(func(name, value interface{}) bool {
		res = append(res, storage.MetricsItemGauage{
			Name:  name.(string),
			Value: value.(float64),
		})
		return true
	})

	return res, nil
}
