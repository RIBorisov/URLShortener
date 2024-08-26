package config

import (
	"flag"
)

var f ServiceConfig

func parseFlags() *ServiceConfig {
	if !flag.Parsed() {
		flag.StringVar(&f.ServerAddress, "a", "localhost:8080", "Address and port to run server, example: localhost:8080")
		flag.StringVar(&f.BaseURL, "b", "http://localhost:8080", "Server address")
		flag.StringVar(&f.FileStoragePath, "f", "", "File path to save data")
		flag.StringVar(&f.DatabaseDSN, "d", "", "Database DSN")
		flag.StringVar(&f.SecretKey, "s", "", "Secret key")
		flag.Parse()
	}
	return &f
}
