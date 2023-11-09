package middleware_test

import (
	"compress/gzip"
	"go-metrics/internal/server/handlers/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareCorrectlySetsContentEncodingHeaderToGzipWhenAcceptEncodingHeaderIncludesGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	middlewareHandler := middleware.CompressMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding header to be 'gzip', but got '%s'", resp.Header.Get("Content-Encoding"))
	}

	body, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello, World!"
	actual := string(data)
	if actual != expected {
		t.Errorf("Expected response body to be '%s', but got '%s'", expected, actual)
	}
}
