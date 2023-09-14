package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"net/http"
	"strconv"
)

func (h MuxHandlers) UpdateCounterMetricsHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	value, err := strconv.ParseInt(vars["value"], 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	it := storage.MetricsItemCounter{
		Name:  vars["name"],
		Value: value,
	}

	err = h.s.AddCounter(it)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.setDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h MuxHandlers) GetCounterMetricsHandle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	item, err := h.s.FindCounterItem(vars["name"])

	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	_, err = w.Write([]byte(strconv.FormatInt(item.Value, 10)))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.setDefaultHeaders(w)
}
