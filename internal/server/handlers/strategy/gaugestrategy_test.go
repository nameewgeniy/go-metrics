package strategy

import (
	"errors"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"github.com/nameewgeniy/go-metrics/internal/server/storage/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddMetric_SuccessGauge(t *testing.T) {
	// Arrange
	name := "test"
	value := "10"
	expectedItem := storage.MetricsItemGauge{
		Name:  name,
		Value: 10.0,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddGauge", expectedItem).Return(nil)

	strategy := &GaugeMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(name, value, storageMock)

	// Assert
	assert.NoError(t, err)
	storageMock.AssertCalled(t, "AddGauge", expectedItem)
}

func TestAddMetric_ParseErrorGauge(t *testing.T) {
	// Arrange
	name := "test"
	value := "not_an_float"

	storageMock := &mock.MockStorage{}
	strategy := &GaugeMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(name, value, storageMock)

	// Assert
	assert.Error(t, err)
}

func TestAddMetric_StorageErrorGauge(t *testing.T) {
	// Arrange
	name := "test"
	value := "10.0"
	expectedError := errors.New("storage error")

	expectedItem := storage.MetricsItemGauge{
		Name:  name,
		Value: 10.0,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddGauge", expectedItem).Return(expectedError)

	strategy := &GaugeMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(name, value, storageMock)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	storageMock.AssertCalled(t, "AddGauge", expectedItem)
}
