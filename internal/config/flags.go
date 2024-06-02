package config

import (
	"flag"
)

var f ServiceConfig

func parseFlags() *ServiceConfig {
	if !flag.Parsed() {
		flag.StringVar(&f.ServerAddress, "a", "localhost:8080", "address and port to run server, example: localhost:8080")
		flag.StringVar(&f.BaseURL, "b", "http://localhost:8080", "server address")
		flag.StringVar(&f.FileStoragePath, "f", "", "file path to save data")
		flag.Parse()
	}
	return &f
}
