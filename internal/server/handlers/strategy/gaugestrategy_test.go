package strategy

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/mock"
	"go-metrics/internal/shared/metrics"
	"testing"
)

func TestAddMetric_SuccessGauge(t *testing.T) {

	// Arrange
	name := "test"
	value := float64(10)

	m := metrics.Metrics{
		ID:    name,
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}

	expectedItem := storage.MetricsItemGauge{
		Name:  name,
		Value: 10.0,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddGauge", expectedItem).Return(nil)

	strategy := &GaugeMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(m, storageMock)

	// Assert
	assert.NoError(t, err)
	storageMock.AssertCalled(t, "AddGauge", expectedItem)
}

func TestAddMetric_StorageErrorGauge(t *testing.T) {
	// Arrange
	name := "test"
	value := float64(10)
	expectedError := errors.New("storage error")

	m := metrics.Metrics{
		ID:    name,
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}

	expectedItem := storage.MetricsItemGauge{
		Name:  name,
		Value: 10.0,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddGauge", expectedItem).Return(expectedError)

	strategy := &GaugeMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(m, storageMock)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	storageMock.AssertCalled(t, "AddGauge", expectedItem)
}
