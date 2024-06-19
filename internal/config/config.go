package config

import (
	"os"
)

const (
	defaultFilePath = "/tmp/short-url-db.json"
	dbDSN           = "DATABASE_DSN"
	baseURL         = "BASE_URL"
	serverAddress   = "SERVER_ADDRESS"
	fileStoragePath = "FILE_STORAGE_PATH"
)

type ServiceConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" env-default:":8080"`
	BaseURL         string `env:"BASE_URL" env-default:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" env-default:"/tmp/short-url-db.json"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

type URLDetail struct {
	Length int `env:"URL_LENGTH" env-default:"8"`
}

type Config struct {
	Service ServiceConfig
	URL     URLDetail
}

func LoadConfig() *Config {
	var cfg Config
	f := parseFlags()
	cfg.URL.Length = 8
	cfg.Service.FileStoragePath = defaultFilePath

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

	return &cfg
}
