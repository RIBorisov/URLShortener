package config

import (
	"flag"
)

var c Config

func parseFlags() *Config {
	if !flag.Parsed() {
		flag.StringVar(&c.App.ServerAddress, "a", ":8080", "Address where server runs")
		flag.StringVar(&c.App.BaseURL, "b", "http://localhost:8080", "Server address")
		flag.StringVar(&c.App.FileStoragePath, "f", "", "File path to save data")
		flag.StringVar(&c.App.DatabaseDSN, "d", "", "Database DSN")
		flag.StringVar(&c.Service.SecretKey, "secret", "", "Secret key")
		flag.BoolVar(&c.App.EnableHTTPS, "s", false, "Enable HTTPS")
		flag.StringVar(&c.App.ConfigFilePath, "c", "", "Config file path")
		flag.StringVar(&c.App.TrustedSubnet, "t", "", "Trusted subnet")
		flag.Parse()
	}
	return &c
}
