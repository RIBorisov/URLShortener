package config

import (
	"os"
)

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" env-default:":8080"`
	BaseURL       string `env:"BASE_URL" env-default:"http://localhost:8080"`
}

type URLDetail struct {
	Length int `env:"URL_LENGTH" env-default:"8"`
}

type Config struct {
	Server ServerConfig
	URL    URLDetail
}

func LoadConfig() *Config {
	var cfg Config

	f := parseFlags()
	cfg.URL.Length = 8

	envBaseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		cfg.Server.BaseURL = envBaseURL
	} else {
		cfg.Server.BaseURL = f.BaseURL
	}

	envAddr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		cfg.Server.ServerAddress = envAddr
	} else {
		cfg.Server.ServerAddress = f.ServerAddress
	}
	return &cfg
}
