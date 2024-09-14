// Package config uses for configure project.
package config

import (
	"os"
	"time"
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
	ServerAddress             string        `env:"SERVER_ADDRESS" env-default:":8080"`
	BaseURL                   string        `env:"BASE_URL" env-default:"http://localhost:8080"`
	FileStoragePath           string        `env:"FILE_STORAGE_PATH" env-default:"/tmp/short-url-db.json"`
	DatabaseDSN               string        `env:"DATABASE_DSN"`
	SecretKey                 string        `env:"SECRET_KEY"`
	BackgroundCleanup         bool          `env:"BACKGROUND_CLEANUP"`
	BackgroundCleanupInterval time.Duration `env:"BACKGROUND_CLEANUP_INTERVAL"`
}

// Config contains main config structures.
type Config struct {
	Service ServiceConfig
}

// LoadConfig loads the config.
func LoadConfig() *Config {
	var cfg Config
	f := parseFlags()
	cfg.Service.FileStoragePath = defaultFilePath
	cfg.Service.SecretKey = secretKeyValue

	sKey, ok := os.LookupEnv(secretKey)
	if ok {
		cfg.Service.SecretKey = sKey
	}

	dsn, ok := os.LookupEnv(dbDSN)
	if ok {
		cfg.Service.DatabaseDSN = dsn
	} else {
		cfg.Service.DatabaseDSN = f.DatabaseDSN
	}

	envBaseURL, ok := os.LookupEnv(baseURL)
	if ok {
		cfg.Service.BaseURL = envBaseURL
	} else {
		cfg.Service.BaseURL = f.BaseURL
	}

	envAddr, ok := os.LookupEnv(serverAddress)
	if ok {
		cfg.Service.ServerAddress = envAddr
	} else {
		cfg.Service.ServerAddress = f.ServerAddress
	}

	path, ok := os.LookupEnv(fileStoragePath)
	if ok {
		cfg.Service.FileStoragePath = path
	} else if f.FileStoragePath != "" {
		cfg.Service.FileStoragePath = f.FileStoragePath
	}

	_, ok = os.LookupEnv(backgroundCleanup)
	if ok {
		cfg.Service.BackgroundCleanup = true
		cfg.Service.BackgroundCleanupInterval = time.Second * 60
	}

	return &cfg
}
