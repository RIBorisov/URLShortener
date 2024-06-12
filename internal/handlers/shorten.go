package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func ShortenHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.ShortenRequest
		ctx := r.Context()
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			logger.Err("failed to decode request body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		short, err := svc.SaveURL(ctx, req.URL)
		if err != nil {
			var duplicateErr *storage.DuplicateRecordError
			if errors.As(err, &duplicateErr) {
				w.WriteHeader(http.StatusConflict)
				short = duplicateErr.Message
			} else {
				logger.Err("failed to save url", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		resultURL, err := url.JoinPath(svc.BaseURL, short)
		if err != nil {
			logger.Err("failed to join path to get result URL", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		resp := models.ShortenResponse{
			Result: resultURL,
		}
		enc := json.NewEncoder(w)
		if err = enc.Encode(resp); err != nil {
			logger.Err("failed to encode response", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
