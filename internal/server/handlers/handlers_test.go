package handlers

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go-metrics/internal/server/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	AddCounterFn      func(i storage.MetricsItemCounter) error
	AddGaugeFn        func(gauge storage.MetricsItemGauge) error
	FindGaugeItemFn   func(name string) (storage.MetricsItemGauge, error)
	FindGaugeAllFn    func() ([]storage.MetricsItemGauge, error)
	FindCounterItemFn func(name string) (storage.MetricsItemCounter, error)
	FindCounterAllFn  func() ([]storage.MetricsItemCounter, error)
	RestoreFn         func() error
	SnapshotFn        func() error
}

func (m *MockStorage) AddCounter(i storage.MetricsItemCounter) error {
	return m.AddCounterFn(i)
}

func (m *MockStorage) AddGauge(gauge storage.MetricsItemGauge) error {
	return m.AddGaugeFn(gauge)
}

func (m *MockStorage) FindGaugeItem(name string) (storage.MetricsItemGauge, error) {
	return m.FindGaugeItemFn(name)
}

func (m *MockStorage) FindGaugeAll() ([]storage.MetricsItemGauge, error) {
	return m.FindGaugeAllFn()
}

func (m *MockStorage) FindCounterItem(name string) (storage.MetricsItemCounter, error) {
	return m.FindCounterItemFn(name)
}

func (m *MockStorage) FindCounterAll() ([]storage.MetricsItemCounter, error) {
	return m.FindCounterAllFn()
}

func (m *MockStorage) Restore() error {
	return m.RestoreFn()
}

func (m *MockStorage) Snapshot() error {
	return m.SnapshotFn()
}

func TestMuxHandlers_UpdateCounterMetricsHandle(t *testing.T) {
	h := MuxHandlers{
		s: &MockStorage{
			AddCounterFn: func(i storage.MetricsItemCounter) error {
				return nil
			},
		},
	}

	// Создаем тестовый HTTP-запрос с нужными параметрами
	req := httptest.NewRequest(http.MethodGet, "/metrics/counter/test/10", nil)

	// Создаем "фейковый" ResponseWriter для записи ответа
	rec := httptest.NewRecorder()

	// Создаем "фейковый" маршрутизатор и регистрируем обработчик
	r := mux.NewRouter()
	r.HandleFunc("/metrics/{type}/{name}/{value}", h.UpdateMetricsHandle)

	// Выполняем запрос с помощью маршрутизатора
	r.ServeHTTP(rec, req)

	// Проверяем код состояния ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedContentType := "text/plain; charset=utf-8"
	if got := rec.Header().Get("content-type"); got != expectedContentType {
		t.Errorf("Expected content type %q, got %q", expectedContentType, got)
	}
}

func TestMuxHandlers_UpdateGaugeMetricsHandle(t *testing.T) {
	h := MuxHandlers{
		s: &MockStorage{
			AddGaugeFn: func(i storage.MetricsItemGauge) error {
				return nil
			},
		},
	}

	// Создаем тестовый HTTP-запрос с нужными параметрами
	req := httptest.NewRequest(http.MethodGet, "/metrics/gauge/test/10", nil)

	// Создаем "фейковый" ResponseWriter для записи ответа
	rec := httptest.NewRecorder()

	// Создаем "фейковый" маршрутизатор и регистрируем обработчик
	r := mux.NewRouter()
	r.HandleFunc("/metrics/{type}/{name}/{value}", h.UpdateMetricsHandle)

	// Выполняем запрос с помощью маршрутизатора
	r.ServeHTTP(rec, req)

	// Проверяем код состояния ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedContentType := "text/plain; charset=utf-8"
	if got := rec.Header().Get("content-type"); got != expectedContentType {
		t.Errorf("Expected content type %q, got %q", expectedContentType, got)
	}
}
