package config

import (
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
	envFileStorage, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cfg.Service.FileStoragePath = envFileStorage
	} else {
		cfg.Service.FileStoragePath = f.FileStoragePath
	}
	return &cfg
}
