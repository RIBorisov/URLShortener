package handlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"shortener/internal/logger"
	"shortener/internal/service"
)

func SaveHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body io.Reader = r.Body
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// если в запросе пришел заголовок Content-Encoding: gzip, пробуем декомпрессить
			// создаем Reader
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Err("failed to create gzip reader", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			defer gr.Close()

			// создаем приемщик прочитанного декомпресса
			var buf bytes.Buffer
			_, err = io.Copy(&buf, gr)
			if err != nil {
				logger.Err("failed to decompress gzip body", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			body = &buf
		}

		long, err := io.ReadAll(body)
		if err != nil {
			logger.Err("failed to read body", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		short := svc.SaveURL(string(long))

		resultURL, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			log.Printf("failed to join path to get result URL: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resultURL))
		if err != nil {
			log.Printf("failed to write the full URL response to client: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
