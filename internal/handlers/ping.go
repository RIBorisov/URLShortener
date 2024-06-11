package handlers

import (
	"context"
	"net/http"

	"shortener/internal/logger"
	"shortener/internal/service"
)

func PingHandler(ctx context.Context, svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Storage.Ping(ctx); err != nil {
			logger.Err("failed to ping Pool", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
