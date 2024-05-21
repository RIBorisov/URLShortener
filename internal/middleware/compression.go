package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

//
//import (
//	"compress/gzip"
//	"io"
//	"net/http"
//	"strings"
//)
//
//type gzipWriter struct {
//	http.ResponseWriter
//	Writer io.Writer
//}
//
//func (w gzipWriter) Write(b []byte) (int, error) {
//	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
//	return w.Writer.Write(b)
//}
//
//func GzipHandler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
//			next.ServeHTTP(w, r)
//			return
//		}
//		if !shouldCompress(r.Header.Get("Content-Type")) {
//			next.ServeHTTP(w, r)
//			return
//		}
//
//		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
//		if err != nil {
//			io.WriteString(w, err.Error())
//			return
//		}
//		defer gz.Close()
//
//		w.Header().Set("Content-Encoding", "gzip")
//		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
//	})
//}

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
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				slog.Error("failed to compress request body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	})
}
