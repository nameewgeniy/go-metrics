package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/strategy"
	"go-metrics/internal/shared"
	"go-metrics/internal/shared/metrics"
	"io"
	"net/http"
)

func (h MuxHandlers) UpdateMetricsHandle(w http.ResponseWriter, r *http.Request) {

	m, err := metrics.NewMetricsFactory().
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

	w.Header().Set("content-type", "application/json")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := metrics.NewMetricsFactory().
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

	err = h.getMetrics(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := m.MarshalJSON()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(content)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) updateMetrics(metrics metrics.Metrics) error {

	strategies := map[string]strategy.MetricsItemStrategy{
		shared.GaugeType:   &strategy.GaugeMetricsItemStrategy{},
		shared.CounterType: &strategy.CounterMetricsItemStrategy{},
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
