package handlers

import (
	"encoding/json"
	"net/http"

	"shortener/internal/models"
	"shortener/internal/service"
)

func GetURLsHandler(svc *service.Service, user *models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urls, err := svc.GetUserURLs(ctx, user)
		if err != nil {
			svc.Log.Err("failed get user urls: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}
		if err = json.NewEncoder(w).Encode(urls); err != nil {
			svc.Log.Err("failed to encode response: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
