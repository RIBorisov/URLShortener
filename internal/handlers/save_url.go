package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"shortener/internal/service"
)

func SaveHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const length = 8
		long, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read body: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		short := svc.GenerateUniqueShortLink(length)
		svc.SaveURL(short, string(long))

		resultURL, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			log.Printf("failed to join path to get result URL: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resultURL))
		if err != nil {
			log.Printf("failed to write the full URL response to client: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
