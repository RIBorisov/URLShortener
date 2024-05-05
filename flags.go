package shortener

import (
	"flag"
)

var FlagRunAddr string
var FlagRunBaseAddr string
var flagsParsed bool

func ParseFlags() {
	if !flagsParsed {
		flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server, example: localhost:8080")
		flag.StringVar(&FlagRunBaseAddr, "b", "http://localhost:8080", "server address")
		flag.Parse()
		flagsParsed = true
	}
}
