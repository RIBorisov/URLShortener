package middleware

import (
	"net/http"
	"shortener/internal/logger"
	"strings"
)

func allowedContentTypes(header string) bool {
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
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			if allowedContentTypes(r.Header.Get("Content-Type")) {
				cw := newCompressWriter(w)
				ow = cw
				defer func() {
					err := cw.Close()
					if err != nil {
						logger.Err("failed to close compress writer", err)
						http.Error(w, "", http.StatusInternalServerError)
						return
					}
				}()
			}
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				logger.Err("failed to read compressed body", err)
				//http.Error(w, "", http.StatusInternalServerError)
				//return
			} else {
				r.Body = cr
				defer func() {
					err = cr.Close()
					if err != nil {
						logger.Err("failed to close compress reader", err)
						http.Error(w, "", http.StatusInternalServerError)
						return
					}
				}()
			}
		}
		next.ServeHTTP(ow, r)
	})
}
