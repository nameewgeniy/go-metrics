package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockMetricsUpdater struct {
	UpdateFn func(metricType, metricName, metricValue string) error
}

func (m *MockMetricsUpdater) Update(metricType, metricName, metricValue string) error {
	return m.UpdateFn(metricType, metricName, metricValue)
}

func TestUpdateMetricsHandle(t *testing.T) {
	handler := Handlers{
		m: &MockMetricsUpdater{
			UpdateFn: func(metricType, metricName, metricValue string) error {
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
