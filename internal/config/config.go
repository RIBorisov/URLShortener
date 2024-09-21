// Package config uses for configure project.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ilyakaznacheev/cleanenv"
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
	ConfigFilePath   string ``
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

	fromFile := &Config{}
	fPath, ok := os.LookupEnv("CONFIG")
	if ok {
		fromFile = LoadConfigFromFile(fPath)
	} else {
		if f.App.ConfigFilePath != "" {
			fromFile = LoadConfigFromFile(f.App.ConfigFilePath)
		}
	}

	sKey, ok := os.LookupEnv(secretKey)
	if ok {
		cfg.Service.SecretKey = sKey
	}

	dsn, ok := os.LookupEnv(dbDSN)
	if ok {
		cfg.App.DatabaseDSN = dsn
	} else if f.App.DatabaseDSN != "" {
		cfg.App.DatabaseDSN = f.App.DatabaseDSN
	} else {
		cfg.App.DatabaseDSN = fromFile.App.DatabaseDSN
	}

	envBaseURL, ok := os.LookupEnv(baseURL)
	if ok {
		cfg.App.BaseURL = envBaseURL
	} else if f.App.BaseURL != "" {
		cfg.App.BaseURL = f.App.BaseURL
	} else {
		cfg.App.BaseURL = fromFile.App.BaseURL
	}

	envAddr, ok := os.LookupEnv(serverAddress)
	if ok {
		cfg.App.ServerAddress = envAddr
	} else if f.App.ServerAddress != "" {
		cfg.App.ServerAddress = f.App.ServerAddress
	} else {
		cfg.App.ServerAddress = fromFile.App.ServerAddress
	}

	path, ok := os.LookupEnv(fileStoragePath)
	if ok {
		cfg.App.FileStoragePath = path
	} else if f.App.FileStoragePath != "" {
		cfg.App.FileStoragePath = f.App.FileStoragePath
	} else {
		cfg.App.FileStoragePath = fromFile.App.FileStoragePath
	}

	v, ok := os.LookupEnv("EnableHTTPS")
	if !ok {
		if f.App.EnableHTTPS || fromFile.App.EnableHTTPS {
			cfg.App.EnableHTTPS = true
		}
	} else {
		val, err := strconv.ParseBool(v)
		if err != nil {
			cfg.App.EnableHTTPS = false
		} else {
			cfg.App.EnableHTTPS = val
		}
	}

	_, ok = os.LookupEnv(backgroundCleanup)
	if ok {
		cfg.Service.BackgroundCleanup = true
		cfg.Service.BackgroundCleanupInterval = time.Second * 60
	}

	return cfg
}

func LoadConfigFromFile(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
