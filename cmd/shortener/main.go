// Package main contains the main entry point for the URL shortener application.
package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"time"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
)

// main is the entry point for the URL shortener application.
//
//go:generate go run ./../metadata/gen_meta.go
func main() {
	log := &logger.Log{}
	log.Initialize("INFO")
	ctx := context.Background()
	if err := initApp(ctx, log); err != nil {
		log.Fatal("unexpected error occurred when initializing application", err)
	}
}

// initApp initializes the URL shortener application.
//
// It loads the configuration, initializes the storage, and starts the HTTP server.
func initApp(ctx context.Context, log *logger.Log) error {
	cfg := config.LoadConfig()
	store, err := storage.LoadStorage(ctx, cfg, log)
	if err != nil {
		log.Err("failed to load storage: ", err)
	}
	defer func() {
		if err = store.Close(); err != nil {
			log.Err("failed to close the connection: ", err)
		}
	}()

	svc := &service.Service{
		Storage:         store,
		BaseURL:         cfg.Service.BaseURL,
		FileStoragePath: cfg.Service.FileStoragePath,
		DatabaseDSN:     cfg.Service.DatabaseDSN,
		Log:             log,
		SecretKey:       cfg.Service.SecretKey,
	}

	if cfg.Service.BackgroundCleanup {
		interval := cfg.Service.BackgroundCleanupInterval
		log.Info("starting storage cleanup task", "period", interval)
		runBackgroundCleanupDB(ctx, store, log, interval)
	}

	r := handlers.NewRouter(svc)

	srv := &http.Server{
		Addr:    cfg.Service.ServerAddress,
		Handler: r,
	}
	log.Info(
		"server starting...",
		slog.String("host", cfg.Service.ServerAddress),
		slog.String("base URL", cfg.Service.BaseURL),
		slog.String("version", buildVersion),
		slog.String("date", buildDate),
		slog.String("commit", buildCommit),
	)

	return srv.ListenAndServe()
}

// runBackgroundCleanupDB runs a background task to clean up the storage periodically.
//
// It uses a ticker to schedule the cleanup at the specified interval.
func runBackgroundCleanupDB(ctx context.Context, store service.URLStorage, log *logger.Log, interval time.Duration) {
	const op = "background cleanup task."

	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			urls, err := store.Cleanup(ctx)
			if err != nil {
				log.Err("failed cleanup storage", err)
			}
			if len(urls) > 0 {
				log.Info(op+" The following url IDs were deleted from storage", "URLs", urls)
			} else {
				log.Info(op+" Nothing to delete. Going to sleep", "time", interval)
			}
		}
	}()
}
