package middleware

import (
	"compress/gzip"
	"go-metrics/internal/server/handlers/middleware/internal/responsewriter"
	"io"
	"net/http"
	"strings"
)

// CompressMiddleware — middleware для сжатия запроса.
func CompressMiddleware(next http.HandlerFunc) http.HandlerFunc {

	zipFn := func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		crw := responsewriter.NewCompressResponseWriter(w, gz)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(crw, r)
	}
	return zipFn
}
