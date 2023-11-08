package middleware_test

import (
	"bytes"
	"go-metrics/internal/server/handlers/middleware"
	"go-metrics/internal/shared"
	"go-metrics/internal/shared/signature"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareSuccessfullyValidatesHashSumOfRequestBody(t *testing.T) {

	signature.Singleton("")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	checkSignatureHandler := middleware.CheckSignature(handler)

	body := []byte("test body")
	hash := signature.Sign.Hash(body)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(shared.HashHeaderName, hash)

	rr := httptest.NewRecorder()

	checkSignatureHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
