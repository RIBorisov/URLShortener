package main

import (
	"log/slog"
	"net/http"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func main() {
	log := logger.Initialize()

	cfg := config.LoadConfig()
	db, err := storage.LoadStorage(cfg)
	if err != nil {
		log.Error("failed to load storage", err)
	}
	svc := &service.Service{DB: db, BaseURL: cfg.Service.BaseURL, FileStoragePath: cfg.Service.FileStoragePath}

	r := handlers.NewRouter(svc, log)

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
