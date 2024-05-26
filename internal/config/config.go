package config

import (
	"log"
	"os"
)

type ServiceConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" env-default:":8080"`
	BaseURL         string `env:"BASE_URL" env-default:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" env-default:"/tmp/short-url-db.json"`
}

type URLDetail struct {
	Length int `env:"URL_LENGTH" env-default:"8"`
}

type Config struct {
	Service ServiceConfig
	URL     URLDetail
}

func LoadConfig() *Config {
	const defaultFilePath = "/tmp/short-url-db.json"
	var cfg Config

	f := parseFlags()
	cfg.URL.Length = 8

	envBaseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		cfg.Service.BaseURL = envBaseURL
	} else {
		cfg.Service.BaseURL = f.BaseURL
	}

	envAddr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		cfg.Service.ServerAddress = envAddr
	} else {
		cfg.Service.ServerAddress = f.ServerAddress
	}

	if f.FileStoragePath != "" {
		envFileStorage := os.Getenv("FILE_STORAGE_PATH")
		switch {
		case envFileStorage != "":
			log.Printf(">>> envFileStorage: %s", envFileStorage)
			cfg.Service.FileStoragePath = envFileStorage
		case f.FileStoragePath != "":
			log.Printf(">>> envFileStorage: %s", f.FileStoragePath)
			cfg.Service.FileStoragePath = f.FileStoragePath
		default:
			log.Printf(">>> defaultFilePath: %s", defaultFilePath)
			cfg.Service.FileStoragePath = defaultFilePath
		}
	}

	return &cfg
}
