package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
)

func BatchHandler(ctx context.Context, svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []models.BatchRequest
		// обрабатываем вход
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			logger.Err("failed to decode request body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Err("failed to close request body", err)
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
			logger.Err("failed to save urls", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		err = enc.Encode(saved)
		if err != nil {
			logger.Err("failed to encode response", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
