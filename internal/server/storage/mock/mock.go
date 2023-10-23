package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go-metrics/internal/server/storage"
)

type MockStorage struct {
	mock.Mock
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) AddGauge(item storage.MetricsItemGauge) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockStorage) FindGaugeItem(name string) (storage.MetricsItemGauge, error) {
	args := m.Called(name)
	return args.Get(0).(storage.MetricsItemGauge), args.Error(1)
}

func (m *MockStorage) FindGaugeAll() ([]storage.MetricsItemGauge, error) {
	args := m.Called()
	return args.Get(0).([]storage.MetricsItemGauge), args.Error(1)
}

func (m *MockStorage) AddCounter(item storage.MetricsItemCounter) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockStorage) FindCounterItem(name string) (storage.MetricsItemCounter, error) {
	args := m.Called(name)
	return args.Get(0).(storage.MetricsItemCounter), args.Error(1)
}

func (m *MockStorage) FindCounterAll() ([]storage.MetricsItemCounter, error) {
	args := m.Called()
	return args.Get(0).([]storage.MetricsItemCounter), args.Error(1)
}

func (m *MockStorage) AddBatchCounters(counters []storage.MetricsItemCounter) error {
	args := m.Called(counters)
	return args.Error(0)
}

func (m *MockStorage) AddBatchGauges(gauges []storage.MetricsItemGauge) error {
	args := m.Called(gauges)
	return args.Error(0)
}

func (m *MockStorage) Restore() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Snapshot() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Up(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Down(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}
