package config

import (
	"flag"
)

var flagRunAddr string
var flagRunBaseAddr string
var flagsParsed bool

func parseFlags() {
	if !flagsParsed {
		flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server, example: localhost:8080")
		flag.StringVar(&flagRunBaseAddr, "b", "http://localhost:8080", "server address")
		flag.Parse()
		flagsParsed = true
	}
}
