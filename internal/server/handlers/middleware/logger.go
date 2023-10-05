package middleware

import (
	"go-metrics/internal/server/handlers/middleware/internal/responsewriter"
	"go-metrics/internal/shared/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(next http.HandlerFunc) http.HandlerFunc {

	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := responsewriter.NewLoggerResponseWriter(w)

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logger.Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("duration", duration.String()),
			zap.Int("status", lw.Info.Status()),
			zap.Int("size", lw.Info.Size()),
		)
	}
	return logFn
}
