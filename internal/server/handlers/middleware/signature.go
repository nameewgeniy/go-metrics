package middleware

import (
	"bytes"
	"errors"
	"go-metrics/internal/shared"
	"go-metrics/internal/shared/signature"
	"io"
	"net/http"
)

var ErrHashNotValid = errors.New("hash not valid")

// CheckSignature — middleware проверяет hash сумму тела запроса.
func CheckSignature(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		hash := r.Header.Get(shared.HashHeaderName)

		if hash == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		ok, err := signature.Sign.Valid(hash, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, ErrHashNotValid.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body)) // Восстанавливаем тело запро
		next.ServeHTTP(w, r)
	}
}
