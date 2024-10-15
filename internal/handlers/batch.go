package handlers

import (
	"encoding/json"
	"net/http"

	"shortener/internal/models"
	"shortener/internal/service"
)

// BatchHandler represents a handler for batch URL shortening requests.
// It handles HTTP POST requests to the /api/shorten/batch endpoint,
// decoding the request body into a slice of BatchRequest objects.
// If the request body is empty or cannot be decoded, it returns an appropriate error response.
// Otherwise, it calls the SaveURLs method of the service to save the URLs and returns the saved URLs in JSON format.
//
// Example usage:
//
// ```go
// handler := BatchHandler(svc)
// http.Handle("/api/shorten/batch", handler)
// ```.
func BatchHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []models.BatchRequest

		ctx := r.Context()
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			svc.Log.Err("failed to decode request body: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				svc.Log.Err("failed to close request body: ", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}()
		if len(req) == 0 {
			http.Error(w, "Empty request batch", http.StatusBadRequest)
			return
		}

		saved, err := svc.SaveURLs(ctx, req)
		if err != nil {
			svc.Log.Err("failed to save urls: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json") //nolint:goconst //not sure const better
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		err = enc.Encode(saved)
		if err != nil {
			svc.Log.Err("failed to encode response: ", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
