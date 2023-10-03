package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go-metrics/internal"
	"go-metrics/internal/models"
	"go-metrics/internal/server/handlers/strategy"
	"go-metrics/internal/server/storage"
	"io"
	"net/http"
)

func (h MuxHandlers) GetMetricsHandle(w http.ResponseWriter, r *http.Request) {

	m, err := models.NewMetricsFactory().
		MakeFromMapForGetMetrics(mux.Vars(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.getMetrics(m)

	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	val, err := m.ValueByType()

	_, err = w.Write([]byte(val))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) GetMetricsJSONHandle(w http.ResponseWriter, r *http.Request) {

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := models.NewMetricsFactory().
		MakeFromBytesForGetMetrics(bytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.getMetrics(m)

	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

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

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) getMetrics(m *models.Metrics) error {

	strategies := map[string]strategy.MetricsItemStrategy{
		internal.GaugeType:   &strategy.GaugeMetricsItemStrategy{},
		internal.CounterType: &strategy.CounterMetricsItemStrategy{},
	}

	strtg, ok := strategies[m.MType]
	if !ok {
		return fmt.Errorf("unsupported metrics type: %s", m.MType)
	}

	err := strtg.GetMetric(m, h.s)

	if err != nil {
		return err
	}

	return nil
}
