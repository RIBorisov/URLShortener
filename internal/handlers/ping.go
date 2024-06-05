package handlers

import (
	"context"
	"net/http"

	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage/db"
)

func PingHandler(ctx context.Context, svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pool, err := db.NewDB(ctx, svc.DSN)
		if err != nil {
			logger.Err("failed to get new DB", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer pool.Pool.Close()

		if err = pool.Pool.Ping(ctx); err != nil {
			logger.Err("failed to ping DB", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}
