// Package main contains the main entry point for the URL shortener application.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/logger"
	"shortener/internal/service"
	"shortener/internal/storage"
)

var (
	buildVersion = "N/A"
	buildCommit  = "N/A"
	buildDate    = "N/A"
)

// main is the entry point for the URL shortener application.
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
		BaseURL:         cfg.App.BaseURL,
		FileStoragePath: cfg.App.FileStoragePath,
		DatabaseDSN:     cfg.App.DatabaseDSN,
		Log:             log,
		SecretKey:       cfg.Service.SecretKey,
	}

	if cfg.Service.BackgroundCleanup {
		interval := cfg.Service.BackgroundCleanupInterval
		log.Info("starting storage cleanup task", "period", interval)
		runBackgroundCleanupDB(ctx, store, log, interval)
	}

	r := handlers.NewRouter(svc)

	log.Info(
		"server starting...",
		slog.String("base URL", cfg.App.BaseURL),
		slog.String("Build version", buildVersion),
		slog.String("Build date", buildDate),
		slog.String("Build commit", buildCommit),
	)

	if !cfg.App.EnableHTTPS {
		return http.ListenAndServe(cfg.App.ServerAddress, r)
	}
	log.Info("enabling TLS..")
	cert, key, err := prepareTLS(log)
	if err != nil {
		log.Err("failed to prepare TLS cert and key", err)
		return err
	}

	return http.ListenAndServeTLS(cfg.App.ServerAddressTLS, cert, key, r)
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

func prepareTLS(log *logger.Log) (string, string, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cert: %w", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode certificate: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode private key: %w", err)
	}

	certFilePath := "tls/server.crt"
	keyFilePath := "tls/server.key"

	certFile, err := os.OpenFile(certFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", "", fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = certFile.Close()
		if err != nil {
			log.Fatal("failed to close file", err)
		}
	}()
	if _, err = certFile.Write(certPEM.Bytes()); err != nil {
		return "", "", fmt.Errorf("failed to write cert file: %w", err)
	}

	keyFile, err := os.OpenFile(keyFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", "", fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = keyFile.Close()
		if err != nil {
			log.Fatal("failed to close file", err)
		}
	}()
	if _, err = keyFile.Write(privateKeyPEM.Bytes()); err != nil {
		return "", "", fmt.Errorf("failed to write key file: %w", err)
	}

	return certFilePath, keyFilePath, nil
}
