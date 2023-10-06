package strategy

import (
	"go-metrics/internal/server/storage"
	"go-metrics/internal/shared/metrics"
)

type CounterMetricsItemStrategy struct{}

func (ms *CounterMetricsItemStrategy) AddMetric(m metrics.Metrics, s storage.Storage) error {

	it := storage.MetricsItemCounter{
		Name:  m.ID,
		Value: *m.Delta,
	}

	return s.AddCounter(it)
}

func (ms *CounterMetricsItemStrategy) GetMetric(m *metrics.Metrics, s storage.Storage) error {
	item, err := s.FindCounterItem(m.ID)

	if err != nil {
		return err
	}

	m.Delta = &item.Value

	return nil
}
