package handlers

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	AddFn  func(item storage.MetricsItem) error
	FindFn func(mType, name string) (storage.MetricsItem, error)
	AllFn  func() ([]storage.MetricsItem, error)
}

func (m *MockStorage) Add(i storage.MetricsItem) error {
	return m.AddFn(i)
}

func (m *MockStorage) Find(mType, name string) (storage.MetricsItem, error) {
	return m.FindFn(mType, name)
}

func (m *MockStorage) All() ([]storage.MetricsItem, error) {
	return m.AllFn()
}

func TestMuxHandlers_UpdateMetricsHandle(t *testing.T) {
	handler := MuxHandlers{
		s: &MockStorage{
			AddFn: func(item storage.MetricsItem) error {
				return nil
			},
		},
	}

	req, err := http.NewRequest(http.MethodPost, "/metrics/type/name/value", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.UpdateMetricsHandle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	expectedContentType := "text/plain; charset=utf-8"
	if got := rr.Header().Get("content-type"); got != expectedContentType {
		t.Errorf("Expected content type %q, got %q", expectedContentType, got)
	}
}
