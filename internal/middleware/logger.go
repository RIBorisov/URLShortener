package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Info(
				"incoming request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("duration", time.Since(start).String()),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				log.Info(
					"request completed",
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
				)
			}()
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
