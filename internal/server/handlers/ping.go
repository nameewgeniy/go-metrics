package handlers

import (
	"net/http"
)

func (h MuxHandlers) PingHandle(w http.ResponseWriter, r *http.Request) {

	if h.p != nil {
		err := h.p.Ping()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
