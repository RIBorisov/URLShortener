package config

import (
	"flag"
)

type Flags struct {
	RunAddr     string
	RunBaseAddr string
}

var f Flags // TODO: обсудить на 1:1 реализацию.

func parseFlags() *Flags {
	if !flag.Parsed() {
		flag.StringVar(&f.RunAddr, "a", "localhost:8080", "address and port to run server, example: localhost:8080")
		flag.StringVar(&f.RunBaseAddr, "b", "http://localhost:8080", "server address")
		flag.Parse()
	}
	return &f
}
