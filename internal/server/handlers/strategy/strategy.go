package strategy

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
)

type MetricsItemStrategy interface {
	AddMetric(name, value string, s storage.Storage) error
	GetMetric(name string, s storage.Storage) (string, error)
}
