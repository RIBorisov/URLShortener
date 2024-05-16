package config

import (
	"flag"
)

var f ServerConfig

func parseFlags() *ServerConfig {
	if !flag.Parsed() {
		flag.StringVar(&f.ServerAddress, "a", "localhost:8080", "address and port to run server, example: localhost:8080")
		flag.StringVar(&f.BaseURL, "b", "http://localhost:8080", "server address")
		flag.Parse()
	}
	return &f
}
