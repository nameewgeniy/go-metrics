package strategy

import (
	"go-metrics/internal/server/storage"
	"go-metrics/internal/shared/metrics"
)

type GaugeMetricsItemStrategy struct{}

func (ms *GaugeMetricsItemStrategy) AddMetric(m metrics.Metrics, s storage.Storage) error {

	it := storage.MetricsItemGauge{
		Name:  m.ID,
		Value: *m.Value,
	}

	return s.AddGauge(it)
}

func (ms *GaugeMetricsItemStrategy) AddBatchMetric(m []metrics.Metrics, s storage.Storage) error {

	var metricsItems []storage.MetricsItemGauge
	for i := range m {
		metricsItems = append(metricsItems, storage.MetricsItemGauge{
			Name:  m[i].ID,
			Value: *m[i].Value,
		})
	}

	return s.AddBatchGauges(metricsItems)
}

func (ms *GaugeMetricsItemStrategy) GetMetric(m *metrics.Metrics, s storage.Storage) error {
	item, err := s.FindGaugeItem(m.ID)

	if err != nil {
		return err
	}

	m.Value = &item.Value

	return nil
}
