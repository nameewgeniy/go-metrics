package strategy

import (
	"errors"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"github.com/nameewgeniy/go-metrics/internal/server/storage/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddMetric_Success(t *testing.T) {
	// Arrange
	name := "test"
	value := "10"
	expectedItem := storage.MetricsItemCounter{
		Name:  name,
		Value: 10,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddCounter", expectedItem).Return(nil)

	strategy := &CounterMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(name, value, storageMock)

	// Assert
	assert.NoError(t, err)
	storageMock.AssertCalled(t, "AddCounter", expectedItem)
}

func TestAddMetric_ParseError(t *testing.T) {
	// Arrange
	name := "test"
	value := "not_an_int"

	storageMock := &mock.MockStorage{}
	strategy := &CounterMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(name, value, storageMock)

	// Assert
	assert.Error(t, err)
}

func TestAddMetric_StorageError(t *testing.T) {
	// Arrange
	name := "test"
	value := "10"
	expectedError := errors.New("storage error")

	expectedItem := storage.MetricsItemCounter{
		Name:  name,
		Value: 10,
	}

	storageMock := &mock.MockStorage{}
	storageMock.On("AddCounter", expectedItem).Return(expectedError)

	strategy := &CounterMetricsItemStrategy{}

	// Act
	err := strategy.AddMetric(name, value, storageMock)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	storageMock.AssertCalled(t, "AddCounter", expectedItem)
}
