package handlers

import (
	"encoding/json"
	"net/http"

	"shortener/internal/models"
	"shortener/internal/service"
)

// DeleteURLsHandler represents a handler for delete URL requests.
func DeleteURLsHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req models.DeleteURLs
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svc.Log.Err("failed to decode request body: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := svc.DeleteURLs(ctx, req)
		if err != nil {
			svc.Log.Err("failed to delete URLs: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		err = json.NewEncoder(w).Encode(req)
		if err != nil {
			svc.Log.Err("failed to encode response body: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
