package middleware

import (
	"net/http"
	"shortener/internal/logger"
	"strings"
)

func allowedContentType(header string) bool {
	allowed := map[string]struct{}{
		"text/html":        {},
		"application/json": {},
	}
	_, ok := allowed[header]
	return ok
}

type BaseMW struct {
	Log *logger.Log
}

func Gzip(log *logger.Log) *BaseMW {
	return &BaseMW{
		Log: log,
	}
}

func (ng *BaseMW) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cT := r.Header.Get("Content-Type")
			if allowedContentType(cT) {
				w.Header().Set("Content-Type", cT)
				cw := newCompressWriter(w)
				ow = cw
				defer func() {
					err := cw.Close()
					if err != nil {
						ng.Log.Err("failed to close compress writer: ", err)
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
				ng.Log.Err("failed to read compressed body: ", err)
				http.Error(w, "check if gzip data is valid", http.StatusBadRequest)
				return
			} else {
				r.Body = cr
				defer func() {
					err = cr.Close()
					if err != nil {
						ng.Log.Err("failed to close compress reader: ", err)
						http.Error(w, "", http.StatusInternalServerError)
						return
					}
				}()
			}
		}
		next.ServeHTTP(ow, r)
	})
}
