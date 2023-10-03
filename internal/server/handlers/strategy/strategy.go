package strategy

import (
	"go-metrics/internal/models"
	"go-metrics/internal/server/storage"
)

type MetricsItemStrategy interface {
	AddMetric(m models.Metrics, s storage.Storage) error
	GetMetric(m *models.Metrics, s storage.Storage) error
}
