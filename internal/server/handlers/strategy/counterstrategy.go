package strategy

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"strconv"
)

type CounterMetricsItemStrategy struct{}

func (ms *CounterMetricsItemStrategy) AddMetric(name, value string, s storage.Storage) error {
	numValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	it := storage.MetricsItemCounter{
		Name:  name,
		Value: numValue,
	}

	return s.AddCounter(it)
}

func (ms *CounterMetricsItemStrategy) GetMetric(name string, s storage.Storage) (string, error) {
	item, err := s.FindCounterItem(name)

	if err != nil {
		return "", err
	}

	return strconv.FormatInt(item.Value, 10), nil
}
