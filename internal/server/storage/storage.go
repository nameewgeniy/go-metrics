package storage

import (
	"context"
	"errors"
)

var ErrItemNotFound = errors.New("item not found")

type MetricsItemGauge struct {
	Name  string
	Value float64
}

type MetricsItemCounter struct {
	Name  string
	Value int64
}

type Storage interface {
	AddGauge(MetricsItemGauge) error
	FindGaugeItem(name string) (MetricsItemGauge, error)
	FindGaugeAll() ([]MetricsItemGauge, error)

	AddCounter(MetricsItemCounter) error
	FindCounterItem(name string) (MetricsItemCounter, error)
	FindCounterAll() ([]MetricsItemCounter, error)

	Up(ctx context.Context) error
	Down(ctx context.Context) error
}
