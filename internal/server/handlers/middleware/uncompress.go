package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// UnCompressMiddleware — middleware для сжатия запроса.
func UnCompressMiddleware(next http.HandlerFunc) http.HandlerFunc {

	zipFn := func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
			return
		}

		if len(body) == 0 {
			http.Error(w, "Тело запроса пустое", http.StatusBadRequest)
			return
		}

		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			if err == io.EOF {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Ошибка чтения данных в формате gzip: "+err.Error(), http.StatusInternalServerError)
			return
		}

		defer reader.Close()

		// Создаем новый запрос с распакованным телом
		uncompressedRequest := *r
		uncompressedRequest.Body = io.NopCloser(reader)

		next.ServeHTTP(w, &uncompressedRequest)
	}
	return zipFn
}
