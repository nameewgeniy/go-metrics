package handlers

import (
	"html/template"
	"net/http"
)

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

	w.Header().Set("content-type", "text/plain; charset=utf-8")
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
