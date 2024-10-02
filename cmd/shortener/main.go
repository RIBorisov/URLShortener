// Package main contains the main entry point for the URL shortener application.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
	"shortener/internal/tasks"
)

var (
	buildVersion = "N/A"
	buildCommit  = "N/A"
	buildDate    = "N/A"
)

// main is the entry point for the URL shortener application.
func main() {
	log := &logger.Log{}
	log.Initialize("DEBUG")

	if err := initApp(log); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("server closed")
		} else {
			log.Fatal("unexpected error occurred when initializing application", err)
		}
	}
}

// initApp initializes the URL shortener application.
//
// It loads the configuration, initializes the storage, and starts the HTTP server.
func initApp(log *logger.Log) error {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(ctx)

	cfg := config.LoadConfig()
	store, err := storage.LoadStorage(ctx, cfg, log)
	if err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}
	defer func() {
		if err = store.Close(); err != nil {
			log.Err("failed to close the connection: ", err)
		}
	}()

	svc := &service.Service{
		Storage:         store,
		BaseURL:         cfg.App.BaseURL,
		FileStoragePath: cfg.App.FileStoragePath,
		DatabaseDSN:     cfg.App.DatabaseDSN,
		Log:             log,
		SecretKey:       cfg.Service.SecretKey,
	}

	if cfg.Service.BackgroundCleanup {
		interval := cfg.Service.BackgroundCleanupInterval
		g.Go(func() error {
			tasks.Run(ctx, store, log, interval)
			<-ctx.Done()
			return nil
		})
	}

	r := handlers.NewRouter(svc)

	log.Info("server starting...",
		slog.String("base URL", cfg.App.BaseURL),
		slog.String("Build version", buildVersion),
		slog.String("Build date", buildDate),
		slog.String("Build commit", buildCommit),
	)
	srv := &http.Server{
		Addr:    cfg.App.ServerAddress,
		Handler: r,
	}

	// enabling graceful shutdown
	readyToShutdown := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		gracefulShutdown(ctx, log, sigint, srv)
	}()

	// run the server with (or w/o) TLS
	g.Go(func() error {
		if !cfg.App.EnableHTTPS {
			log.Info("TLS disabled..")
			if err = srv.ListenAndServe(); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					return fmt.Errorf("failed listen and serve: %w", err)
				}
			}
		} else {
			if err = srv.ListenAndServeTLS("tls/server.crt", "tls/server.key"); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					return fmt.Errorf("failed listen and serve: %w", err)
				}
			}
		}
		return nil
	})

	<-ctx.Done()
	log.Info("received signal to stop application")

	if err = g.Wait(); err != nil {
		log.Err("failed to wait for all goroutines finished", err)
	}
	log.Info("closed all goroutines, now we may shutdown the server")
	close(readyToShutdown)

	return nil
}

func gracefulShutdown(ctx context.Context, log *logger.Log, sigint chan os.Signal, srv *http.Server) {
	<-sigint
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		log.Err("failed to shutdown server", shutdownErr)
	}
	log.Debug("graceful shutdown complete..")
}
