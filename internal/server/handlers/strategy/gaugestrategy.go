package strategy

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"strconv"
)

type GaugeMetricsItemStrategy struct{}

func (ms *GaugeMetricsItemStrategy) AddMetric(name, value string, s storage.Storage) error {
	numValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	it := storage.MetricsItemGauge{
		Name:  name,
		Value: numValue,
	}

	return s.AddGauge(it)
}

func (ms *GaugeMetricsItemStrategy) GetMetric(name string, s storage.Storage) (string, error) {
	item, err := s.FindGaugeItem(name)

	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(item.Value, 'f', -1, 64), nil
}
