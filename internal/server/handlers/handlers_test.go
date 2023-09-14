package handlers

import (
	"github.com/gorilla/mux"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	AddCounterFn      func(i storage.MetricsItemCounter) error
	AddGauageFn       func(gauage storage.MetricsItemGauage) error
	FindGauageItemFn  func(name string) (storage.MetricsItemGauage, error)
	FindGauageAllFn   func() ([]storage.MetricsItemGauage, error)
	FindCounterItemFn func(name string) (storage.MetricsItemCounter, error)
	FindCounterAllFn  func() ([]storage.MetricsItemCounter, error)
}

func (m *MockStorage) AddCounter(i storage.MetricsItemCounter) error {
	return m.AddCounterFn(i)
}

func (m *MockStorage) AddGauage(gauage storage.MetricsItemGauage) error {
	return m.AddGauageFn(gauage)
}

func (m *MockStorage) FindGauageItem(name string) (storage.MetricsItemGauage, error) {
	return m.FindGauageItemFn(name)
}

func (m *MockStorage) FindGauageAll() ([]storage.MetricsItemGauage, error) {
	return m.FindGauageAllFn()
}

func (m *MockStorage) FindCounterItem(name string) (storage.MetricsItemCounter, error) {
	return m.FindCounterItemFn(name)
}

func (m *MockStorage) FindCounterAll() ([]storage.MetricsItemCounter, error) {
	return m.FindCounterAllFn()
}

func TestMuxHandlers_UpdateCounterMetricsHandle(t *testing.T) {
	h := MuxHandlers{
		s: &MockStorage{
			AddCounterFn: func(i storage.MetricsItemCounter) error {
				return nil
			},
		}, // замените на вашу реализацию хранилища
	}

	// Создаем тестовый HTTP-запрос с нужными параметрами
	req := httptest.NewRequest(http.MethodGet, "/metrics/counter1/10", nil)

	// Создаем "фейковый" ResponseWriter для записи ответа
	rec := httptest.NewRecorder()

	// Создаем "фейковый" маршрутизатор и регистрируем обработчик
	r := mux.NewRouter()
	r.HandleFunc("/metrics/{name}/{value}", h.UpdateCounterMetricsHandle)

	// Выполняем запрос с помощью маршрутизатора
	r.ServeHTTP(rec, req)

	// Проверяем код состояния ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedContentType := "text/plain; charset=utf-8"
	if got := rec.Header().Get("content-type"); got != expectedContentType {
		t.Errorf("Expected content type %q, got %q", expectedContentType, got)
	}
}

func TestMuxHandlers_UpdateGauageMetricsHandle(t *testing.T) {
	h := MuxHandlers{
		s: &MockStorage{
			AddGauageFn: func(i storage.MetricsItemGauage) error {
				return nil
			},
		}, // замените на вашу реализацию хранилища
	}

	// Создаем тестовый HTTP-запрос с нужными параметрами
	req := httptest.NewRequest(http.MethodGet, "/metrics/gauage/10", nil)

	// Создаем "фейковый" ResponseWriter для записи ответа
	rec := httptest.NewRecorder()

	// Создаем "фейковый" маршрутизатор и регистрируем обработчик
	r := mux.NewRouter()
	r.HandleFunc("/metrics/{name}/{value}", h.UpdateGauageMetricsHandle)

	// Выполняем запрос с помощью маршрутизатора
	r.ServeHTTP(rec, req)

	// Проверяем код состояния ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedContentType := "text/plain; charset=utf-8"
	if got := rec.Header().Get("content-type"); got != expectedContentType {
		t.Errorf("Expected content type %q, got %q", expectedContentType, got)
	}
}
