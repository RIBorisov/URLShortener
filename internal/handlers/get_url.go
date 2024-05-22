package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"shortener/internal/service"
)

func GetHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "id")
		long, err := svc.GetURL(short)
		if err != nil {
			log.Printf("failed to get URL: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		origin, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			log.Printf("failed to join path to get redirect URL: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", origin)
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
	}
}
