package handlers

import (
	"net/http"

	"shortener/internal/logger"
	"shortener/internal/service"
)

func PingHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := svc.Storage.Ping(ctx); err != nil {
			logger.Err("failed to ping Pool", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
