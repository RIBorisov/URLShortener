package handlers

import (
	"net/http"
	"net/url"
	"shortener/internal/logger"

	"github.com/go-chi/chi/v5"

	"shortener/internal/service"
)

func GetHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		short := chi.URLParam(r, "id")
		long, err := svc.GetURL(ctx, short)
		if err != nil {
			logger.Err("failed to get URL", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		origin, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			logger.Err("failed to join path to get redirect URL", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", origin)
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
	}
}
