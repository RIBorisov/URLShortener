package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func SaveHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		long, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Err("failed to read body", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		short, err := svc.SaveURL(string(long))
		if err != nil {
			var duplicateErr *storage.DuplicateRecordError
			if errors.As(err, &duplicateErr) {
				w.WriteHeader(http.StatusConflict)
			} else {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
		}

		resultURL, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			log.Printf("failed to join path to get result URL: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(resultURL))
		if err != nil {
			log.Printf("failed to write the full URL response to client: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
