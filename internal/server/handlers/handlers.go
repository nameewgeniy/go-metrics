package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nameewgeniy/go-metrics/internal"
	strategy "github.com/nameewgeniy/go-metrics/internal/server/handlers/strategy"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"html/template"
	"net/http"
)

type MuxHandlers struct {
	s storage.Storage
}

func NewMuxHandlers(s storage.Storage) *MuxHandlers {
	return &MuxHandlers{
		s: s,
	}
}

func (h MuxHandlers) UpdateMetricsHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	strategies := map[string]strategy.MetricsItemStrategy{
		internal.GaugeType:   &strategy.GaugeMetricsItemStrategy{},
		internal.CounterType: &strategy.CounterMetricsItemStrategy{},
	}

	strtg, ok := strategies[vars["type"]]
	if !ok {
		http.Error(w, fmt.Sprintf("unsupported metrics type: %s", vars["type"]), http.StatusBadRequest)
		return
	}

	err := strtg.AddMetric(vars["name"], vars["value"], h.s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.setDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) GetMetricsHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	strategies := map[string]strategy.MetricsItemStrategy{
		internal.GaugeType:   &strategy.GaugeMetricsItemStrategy{},
		internal.CounterType: &strategy.CounterMetricsItemStrategy{},
	}

	strtg, ok := strategies[vars["type"]]
	if !ok {
		http.Error(w, fmt.Sprintf("unsupported metrics type: %s", vars["type"]), http.StatusBadRequest)
		return
	}

	val, err := strtg.GetMetric(vars["name"], h.s)

	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	_, err = w.Write([]byte(val))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.setDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) ViewMetricsHandle(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := h.makeMetricsTemplateData()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон и передаем данные
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.setDefaultHeaders(w)
}

func (h MuxHandlers) makeMetricsTemplateData() (map[string]any, error) {
	data := make(map[string]any)

	counters, err := h.s.FindCounterAll()

	if err != nil {
		return nil, err
	}

	for _, c := range counters {
		data[c.Name] = c.Value
	}

	gauges, err := h.s.FindGaugeAll()

	if err != nil {
		return nil, err
	}

	for _, g := range gauges {
		data[g.Name] = g.Value
	}

	return data, nil
}

func (h MuxHandlers) setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
}
