package handlers

import (
	"net/http"
	"strings"
)

type Metrics interface {
	Update(mType, mName, mValue string) error
}

type Handlers struct {
	m Metrics
}

func NewHandlers(m Metrics) *Handlers {
	return &Handlers{
		m: m,
	}
}

func (h Handlers) UpdateMetricsHande(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := parts[2], parts[3], parts[4]

	err := h.m.Update(metricType, metricName, metricValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
