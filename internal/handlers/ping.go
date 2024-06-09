package handlers

import (
	"context"
	"net/http"

	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage/postgres"
)

func PingHandler(ctx context.Context, svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pool, err := postgres.New(ctx, svc.DatabaseDSN)
		if err != nil {
			logger.Err("failed to get new Pool", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer pool.Pool.Close()

		if err = pool.Pool.Ping(ctx); err != nil {
			logger.Err("failed to ping Pool", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}
