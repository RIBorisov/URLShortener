package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			fmt.Printf("\n>>>r.Headers %+v<<<", r.Header.Clone())
			fmt.Printf("\n>>>r.URL %+v<<<", r.URL)
			fmt.Printf("\n>>>TypeOf(r.URL.Path) %+v<<<", reflect.TypeOf(r.URL.Path))
			fmt.Printf("\n>>>r.URL.Path %+v<<<\n", r.URL.Path)

			defer func() {
				log.Info(
					"OK",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
					slog.String("duration", time.Since(start).String()),
				)
			}()
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
