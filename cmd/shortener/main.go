package main

import (
	"context"
	"log/slog"
	"net/http"

	//_ "github.com/jackc/pgx/v5"

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
	store, err := storage.NewStorage(cfg)
	if err != nil {
		log.Error("failed to load store", err)
	}
	if err != nil {
		log.Error("failed to init DB", err)
	}
	svc := &service.Service{
		Storage:         store,
		BaseURL:         cfg.Service.BaseURL,
		FileStoragePath: cfg.Service.FileStoragePath,
		DSN:             cfg.Service.DatabaseDSN,
	}

	r := handlers.NewRouter(ctx, svc)

	srv := &http.Server{
		Addr:    cfg.Service.ServerAddress,
		Handler: r,
	}

	log.Info(
		"server starting...",
		slog.String("host", cfg.Service.ServerAddress),
		slog.String("baseURL", cfg.Service.BaseURL),
		slog.String("fileStoragePath", cfg.Service.FileStoragePath),
	)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server: %v", err)
	}
}
