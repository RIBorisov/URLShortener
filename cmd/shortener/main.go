package main

import (
	"context"
	"log/slog"
	"net/http"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func main() {
	ctx := context.Background()
	log := logger.Initialize()

	cfg := config.LoadConfig()
	store, err := storage.LoadStorage(ctx, cfg)
	if err != nil {
		logger.Err("failed to load storage", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			logger.Err("failed to close the connection", err)
		}
	}()
	svc := &service.Service{
		Storage:         store,
		BaseURL:         cfg.Service.BaseURL,
		FileStoragePath: cfg.Service.FileStoragePath,
		DatabaseDSN:     cfg.Service.DatabaseDSN,
	}

	r := handlers.NewRouter(svc)

	srv := &http.Server{
		Addr:    cfg.Service.ServerAddress,
		Handler: r,
	}

	log.Info(
		"server starting...",
		slog.String("host", cfg.Service.ServerAddress),
		slog.String("baseURL", cfg.Service.BaseURL),
	)
	if err = srv.ListenAndServe(); err != nil {
		logger.Err("failed to start server", err)
	}
}
