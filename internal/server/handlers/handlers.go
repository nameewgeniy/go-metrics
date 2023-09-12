package handlers

import (
	"errors"
	"github.com/gorilla/mux"
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
	it := storage.MetricsItem{
		Type:  vars["type"],
		Name:  vars["name"],
		Value: vars["value"],
	}

	err := h.s.Add(it)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) ViewMetricsHandle(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := h.s.All()

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

	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) GetMetricsHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	item, err := h.s.Find(vars["type"], vars["name"])

	if err != nil {
		if errors.Is(err, storage.ItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	_, err = w.Write([]byte(item.String()))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
