package middleware_test

import (
	"bytes"
	"compress/gzip"
	"go-metrics/internal/server/handlers/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareDecompressesGzipRequestBodies(t *testing.T) {
	// Create a new request with a gzip-encoded body
	requestBody := "gzip-encoded-body"
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write([]byte(requestBody))
	gz.Close()
	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock handler function
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request body is correctly decompressed
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		if string(body) != requestBody {
			t.Errorf("Expected decompressed body to be %s, but got %s", requestBody, string(body))
		}
		w.WriteHeader(http.StatusOK)
	}

	// Invoke the middleware with the mock handler function
	middleware := middleware.UnCompressMiddleware(handlerFunc)
	middleware.ServeHTTP(rr, req)

	// Verify the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
}
