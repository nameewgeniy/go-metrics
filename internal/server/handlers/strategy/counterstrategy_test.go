package strategy

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go-metrics/internal/models"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/mock"
	"testing"
)

func TestAddMetric_Success(t *testing.T) {

	name := "test"
	delta := int64(10)

	m := models.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
		Value: nil,
	}

	expectedItem := storage.MetricsItemCounter{
		Name:  name,
		Value: delta,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddCounter", expectedItem).Return(nil)

	strategy := &CounterMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(m, storageMock)

	// Assert
	assert.NoError(t, err)
	storageMock.AssertCalled(t, "AddCounter", expectedItem)
}

func TestAddMetric_StorageError(t *testing.T) {
	// Arrange
	name := "test"
	delta := int64(10)

	m := models.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &delta,
		Value: nil,
	}

	expectedError := errors.New("storage error")

	expectedItem := storage.MetricsItemCounter{
		Name:  name,
		Value: 10,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddCounter", expectedItem).Return(expectedError)

	strategy := &CounterMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(m, storageMock)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	storageMock.AssertCalled(t, "AddCounter", expectedItem)
}
