package handlers

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMuxHandlers_UpdateCounterMetricsHandle(t *testing.T) {

	storageMock := mock.NewMockStorage()
	storageMock.On("AddCounter", storage.MetricsItemCounter{
		Name:  "test",
		Value: 10,
	}).Return(nil)

	h := MuxHandlers{
		s: storageMock,
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

	storageMock := mock.NewMockStorage()
	storageMock.On("AddGauge", storage.MetricsItemGauge{
		Name:  "test",
		Value: 10.0,
	}).Return(nil)

	h := MuxHandlers{
		s: storageMock,
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
