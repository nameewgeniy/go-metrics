package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal"
	"go-metrics/internal/models"
	"go-metrics/internal/server/handlers/strategy"
	"io"
	"net/http"
)

func (h MuxHandlers) UpdateMetricsHandle(w http.ResponseWriter, r *http.Request) {

	m, err := models.NewMetricsFactory().
		MakeFromMapForUpdateMetrics(mux.Vars(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.updateMetrics(*m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) UpdateMetricsJSONHandle(w http.ResponseWriter, r *http.Request) {

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := models.NewMetricsFactory().
		MakeFromBytesForUpdateMetrics(bytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.updateMetrics(*m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) updateMetrics(metrics models.Metrics) error {

	strategies := map[string]strategy.MetricsItemStrategy{
		internal.GaugeType:   &strategy.GaugeMetricsItemStrategy{},
		internal.CounterType: &strategy.CounterMetricsItemStrategy{},
	}

	chosenStrategy, ok := strategies[metrics.MType]
	if !ok {
		return fmt.Errorf("unsupported metrics type: %s", metrics.MType)
	}

	err := chosenStrategy.AddMetric(metrics, h.s)
	if err != nil {
		return err
	}

	return nil
}
