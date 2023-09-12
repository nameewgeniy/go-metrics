package goriila

import (
	"github.com/gorilla/mux"
	"github.com/nameewgeniy/go-metrics/internal/server/handlers"
	"net/http"
)

type MuxHandlers struct {
	m handlers.Metrics
}

func NewMuxHandlers(m handlers.Metrics) *MuxHandlers {
	return &MuxHandlers{
		m: m,
	}
}

func (h MuxHandlers) UpdateMetricsHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	metricType, metricName, metricValue := vars["type"], vars["name"], vars["value"]

	err := h.m.Update(metricType, metricName, metricValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
