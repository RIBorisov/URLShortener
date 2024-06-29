package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"shortener/internal/service"
	"shortener/internal/storage"
)

func GetHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		short := chi.URLParam(r, "id")
		long, err := svc.GetURL(ctx, short)
		if err != nil {
			if errors.Is(err, storage.ErrURLDeleted) {
				svc.Log.Info("requested deleted url", "short", short)
				w.WriteHeader(http.StatusGone)
			} else {
				svc.Log.Err("failed to get URL: ", err)
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		origin, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			svc.Log.Err("failed to join path to get redirect URL: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", origin)
		http.Redirect(w, r, long, http.StatusTemporaryRedirect)
	}
}
