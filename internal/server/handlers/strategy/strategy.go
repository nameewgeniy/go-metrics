package strategy

import (
	"go-metrics/internal/server/storage"
	"go-metrics/internal/shared/metrics"
)

type MetricsItemStrategy interface {
	AddMetric(m metrics.Metrics, s storage.Storage) error
	GetMetric(m *metrics.Metrics, s storage.Storage) error
}
