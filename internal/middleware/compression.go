package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

func shouldCompress(header string) bool {
	contentTypes := []string{
		"text/html",
		"application/json",
	}
	for _, item := range contentTypes {
		if item == header {
			return true
		}
	}
	return false
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldCompress(w.Header().Get("Content-Type")) {
			next.ServeHTTP(w, r)
			return
		}
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer func(cw *compressWriter) {
				err := cw.Close()
				if err != nil {
					slog.Error("failed to close compressWriter")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}(cw)
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				slog.Error("failed to read compressed body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func(cr *compressReader) {
				err := cr.Close()
				if err != nil {
					slog.Error("failed to close compressReader")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}(cr)
		}
		next.ServeHTTP(ow, r)
	})
}
