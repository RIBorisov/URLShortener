package config

import (
	"os"
	s "shortener"
)

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" env-default:":8080"`
	BaseURL       string `env:"BASE_URL" env-default:"http://localhost:8080"`
}

type Config struct {
	Server ServerConfig
}

func LoadConfig() *Config {
	var cfg Config

	s.ParseFlags()

	cfg.Server.BaseURL = "http://localhost:8080"
	cfg.Server.ServerAddress = "localhost:8080"

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.Server.BaseURL = envBaseURL
	} else if s.FlagRunBaseAddr != "" {
		cfg.Server.BaseURL = s.FlagRunBaseAddr
	}

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		cfg.Server.ServerAddress = envAddr
	} else if s.FlagRunAddr != ":8080" {
		cfg.Server.ServerAddress = s.FlagRunAddr
	}

	return &cfg
}

//package config
//
//import (
//	"os"
//	s "shortener"
//)
//
//var Config struct {
//	ServerAddress string `env:"SERVER_ADDRESS" env-default:":8080"`
//	BaseURL       string `env:"BASE_URL" env-default:"http://localhost:8080"`
//}
//
//func LoadConfig() {
//	s.ParseFlags()
//	Config.BaseURL = "http://localhost:8080"
//	Config.ServerAddress = "localhost:8080"
//
//	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
//		Config.BaseURL = envBaseURL
//	} else if s.FlagRunBaseAddr != "" {
//		Config.BaseURL = s.FlagRunBaseAddr
//	}
//
//	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
//		Config.ServerAddress = envAddr
//	} else if s.FlagRunAddr != ":8080" {
//		Config.ServerAddress = s.FlagRunAddr
//	}
//}
