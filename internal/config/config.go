// Package config uses for configure project.
package config

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	defaultFilePath   = "/tmp/short-url-db.json"
	dbDSN             = "DATABASE_DSN"
	baseURL           = "BASE_URL"
	serverAddress     = "SERVER_ADDRESS"
	fileStoragePath   = "FILE_STORAGE_PATH"
	secretKey         = "SECRET_KEY"
	secretKeyValue    = "!@#$YdBg0DS"
	backgroundCleanup = "BACKGROUND_CLEANUP"
)

// ServiceConfig contains common config entities.
type ServiceConfig struct {
	SecretKey                 string        `env:"SECRET_KEY"`
	BackgroundCleanup         bool          `env:"BACKGROUND_CLEANUP"`
	BackgroundCleanupInterval time.Duration `env:"BACKGROUND_CLEANUP_INTERVAL"`
}

// AppConfig contains application envs.
type AppConfig struct {
	ServerAddress    string `env:"SERVER_ADDRESS" envDefault:":8080"`
	ServerAddressTLS string `env:"SERVER_ADDRESS_TLS" envDefault:":8443"`
	BaseURL          string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
	DatabaseDSN      string `env:"DATABASE_DSN"`
	EnableHTTPS      bool   `env:"ENABLE_HTTPS" envDefault:"0"`
}

// Config contains main config structures.
type Config struct {
	App     AppConfig
	Service ServiceConfig
}

// LoadConfig loads the config.
func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("failed to parse config")
	}
	f := parseFlags()
	cfg.App.FileStoragePath = defaultFilePath
	cfg.Service.SecretKey = secretKeyValue

	sKey, ok := os.LookupEnv(secretKey)
	if ok {
		cfg.Service.SecretKey = sKey
	}

	dsn, ok := os.LookupEnv(dbDSN)
	if ok {
		cfg.App.DatabaseDSN = dsn
	} else {
		cfg.App.DatabaseDSN = f.App.DatabaseDSN
	}

	envBaseURL, ok := os.LookupEnv(baseURL)
	if ok {
		cfg.App.BaseURL = envBaseURL
	} else {
		cfg.App.BaseURL = f.App.BaseURL
	}

	envAddr, ok := os.LookupEnv(serverAddress)
	if ok {
		cfg.App.ServerAddress = envAddr
	} else {
		cfg.App.ServerAddress = f.App.ServerAddress
	}

	path, ok := os.LookupEnv(fileStoragePath)
	if ok {
		cfg.App.FileStoragePath = path
	} else if f.App.FileStoragePath != "" {
		cfg.App.FileStoragePath = f.App.FileStoragePath
	}

	_, ok = os.LookupEnv(backgroundCleanup)
	if ok {
		cfg.Service.BackgroundCleanup = true
		cfg.Service.BackgroundCleanupInterval = time.Second * 60
	}

	return cfg
}
