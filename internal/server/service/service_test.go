package service

import (
	"errors"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"strconv"
	"testing"
)

type MockStorage struct {
	SaveFn func(item storage.MetricsItem) error
}

func (s *MockStorage) Save(item storage.MetricsItem) error {
	return s.SaveFn(item)
}

func TestUpdate(t *testing.T) {
	mockStorage := &MockStorage{
		SaveFn: func(item storage.MetricsItem) error {
			return nil
		},
	}

	metrics := Metrics{
		s: mockStorage,
	}

	t.Run("Update gauge metric", func(t *testing.T) {
		err := metrics.Update("gauge", "metric1", "2.5")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Update counter metric", func(t *testing.T) {
		err := metrics.Update("counter", "metric2", "10")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Invalid value for gauge metric", func(t *testing.T) {
		err := metrics.Update("gauge", "metric4", "invalid")
		if err == nil {
			t.Error("Expected an error, got nil")
		} else {
			var expectedErr *strconv.NumError
			if !errors.As(err, &expectedErr) {
				t.Errorf("Expected %T error, got %T", expectedErr, err)
			}
		}
	})

	t.Run("Invalid value for counter metric", func(t *testing.T) {
		err := metrics.Update("counter", "metric5", "invalid")
		if err == nil {
			t.Error("Expected an error, got nil")
		} else {
			var expectedErr *strconv.NumError
			if !errors.As(err, &expectedErr) {
				t.Errorf("Expected %T error, got %T", expectedErr, err)
			}
		}
	})

	t.Run("Storage error", func(t *testing.T) {
		mockStorage := &MockStorage{
			SaveFn: func(item storage.MetricsItem) error {
				return errors.New("storage error")
			},
		}

		metrics := Metrics{
			s: mockStorage,
		}

		err := metrics.Update("gauge", "metric6", "3.7")
		if err == nil {
			t.Error("Expected an error, got nil")
		} else {
			if err.Error() != "storage error" {
				t.Errorf("Expected error message %q, got %q", "storage error", err.Error())
			}
		}
	})
}
