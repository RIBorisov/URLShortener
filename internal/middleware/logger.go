package middleware

import (
	"net/http"
	"shortener/internal/logger"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type BaseLog struct {
	Log *logger.Log
}

func Log(log *logger.Log) *BaseLog {
	return &BaseLog{
		Log: log,
	}
}

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
