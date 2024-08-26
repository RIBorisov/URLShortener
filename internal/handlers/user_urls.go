package handlers

import (
	"encoding/json"
	"net/http"

	"shortener/internal/service"
)

// GetURLsHandler represents a handler for getting user URLs requests.
func GetURLsHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urls, err := svc.GetUserURLs(ctx)
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
