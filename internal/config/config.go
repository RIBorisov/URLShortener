package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Println("env CONFIG_PATH is not set, getting a default value")
		configPath = getDefaultPath()
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func getDefaultPath() string {
	// в тестах гитхаб путь формируется неправильно
	//filePath, _ := filepath.Abs(filepath.Dir("../../internal/config/cfg.yaml"))
	//defaultPath := filePath + "/cfg.yaml"
	projectDir, _ := os.Getwd()
	defaultPath := projectDir + "/internal/config/cfg.yaml"
	return defaultPath
}
