package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal/server/handlers/strategy"
	"go-metrics/internal/shared"
	"go-metrics/internal/shared/metrics"
	"io"
	"log"
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
		log.Print(err)
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

func (h MuxHandlers) UpdateBatchMetricsHandle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(string(bytes))

	requestMetrics, err := metrics.NewMetricsFactory().
		MakeFromBytesForBatchUpdateMetrics(bytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var gaugeMetrics []metrics.Metrics
	var counterMetrics []metrics.Metrics

	for i, v := range requestMetrics {
		switch v.MType {
		case shared.GaugeType:
			gaugeMetrics = append(gaugeMetrics, requestMetrics[i])
		case shared.CounterType:
			counterMetrics = append(counterMetrics, requestMetrics[i])
		}
	}

	if len(gaugeMetrics) > 0 {
		gaugeStrategy := strategy.GaugeMetricsItemStrategy{}

		if err = gaugeStrategy.AddBatchMetric(gaugeMetrics, h.s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if len(counterMetrics) > 0 {
		counterStrategy := strategy.CounterMetricsItemStrategy{}

		if err = counterStrategy.AddBatchMetric(counterMetrics, h.s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	mtr, err := h.allMetrics()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := json.Marshal(mtr)

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

func (h MuxHandlers) allMetrics() ([]metrics.Metrics, error) {
	var mtr []metrics.Metrics

	counters, err := h.s.FindCounterAll()

	if err != nil {
		return nil, err
	}

	for _, c := range counters {
		mtr = append(mtr, metrics.Metrics{
			ID:    c.Name,
			MType: shared.CounterType,
			Delta: &c.Value,
		})
	}

	gauges, err := h.s.FindGaugeAll()

	if err != nil {
		return nil, err
	}

	for _, g := range gauges {
		mtr = append(mtr, metrics.Metrics{
			ID:    g.Name,
			MType: shared.GaugeType,
			Value: &g.Value,
		})
	}

	return mtr, nil
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
