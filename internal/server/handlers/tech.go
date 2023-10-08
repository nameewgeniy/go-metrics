package handlers

import (
	"net/http"
)

func (h MuxHandlers) PingDBHandle(w http.ResponseWriter, r *http.Request) {

	err := h.p.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
