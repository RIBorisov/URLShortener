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
	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	log := logger.Initialize()
	svc := &service.Service{DB: db, BaseURL: cfg.Server.BaseURL}

	r := handlers.NewRouter(svc, log)

	srv := &http.Server{
		Addr:    cfg.Server.ServerAddress,
		Handler: r,
	}

	log.Info("server starting ", slog.String("host", cfg.Server.ServerAddress))
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server: %v", err)
	}
}
