package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func main() {
	log := &logger.Log{}
	log.Initialize("INFO")
	ctx := context.Background()
	if err := initApp(ctx, log); err != nil {
		log.Fatal("unexpected error occurred when initializing application", err)
	}
}

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
	stopCh := make(chan struct{})

	if cfg.Service.BackgroundCleanup {
		interval := cfg.Service.BackgroundCleanupInterval
		log.Info("starting storage cleanup task", "period", interval)
		stopCh = runBackgroundCleanupDB(ctx, store, log, interval)
	}

	r := handlers.NewRouter(svc)

	srv := &http.Server{
		Addr:    cfg.Service.ServerAddress,
		Handler: r,
	}
	log.Info(
		"server starting...",
		slog.String("host", cfg.Service.ServerAddress),
		slog.String("BaseURL", cfg.Service.BaseURL),
	)

	// Создаем канал для передачи сигналов об остановке
	sigs := make(chan os.Signal, 1)
	// Регистрируем обработчик сигналов
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Info("Received signal", "signal", sig.String())
		stopCh <- struct{}{}
	}()

	return srv.ListenAndServe()
}

func runBackgroundCleanupDB(
	ctx context.Context, store storage.URLStorage, log *logger.Log, interval time.Duration,
) chan struct{} {
	const op = "background cleanup task."
	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				urls, err := store.Cleanup(ctx)
				if err != nil {
					log.Err("failed cleanup storage", err)
				}
				if len(urls) > 0 {
					log.Info(op+" The following url IDs were deleted from storage", "URLs", urls)
				} else {
					log.Info(op+" Nothing to delete. Going to sleep", "time", interval)
				}
			case <-stopCh:
				log.Info("Stopping " + op)
				ticker.Stop()
				return
			}
		}
	}()

	return stopCh
}
