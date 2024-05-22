package handlers

import (
	"log/slog"
	"net/http"
	"net/url"
	"shortener/internal/logger"

	"github.com/go-chi/chi/v5"

	"shortener/internal/service"
)

func GetHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "id")
		long, err := svc.GetURL(short)
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
		slog.Info(">>> GOT: ", "ORIGIN", origin)
		w.Header().Set("Location", string(origin))
		http.Redirect(w, r, string(long), http.StatusTemporaryRedirect)
	}
}
