package middleware

import (
	"net/http"
	"time"

	"shortener/internal/logger"

	"github.com/go-chi/chi/v5/middleware"
)

// BaseLog represents the base logging middleware.
type BaseLog struct {
	Log *logger.Log
}

// Log creates a new instance of the BaseLog middleware.
func Log(log *logger.Log) *BaseLog {
	return &BaseLog{
		Log: log,
	}
}

// Middleware returns an HTTP handler that logs request details.
func (b *BaseLog) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			b.Log.Info(
				"OK",
				"path", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"size", ww.BytesWritten(),
				"duration", time.Since(start).String(),
			)
		}()
		next.ServeHTTP(ww, r)
	},
	)
}
