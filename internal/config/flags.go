package config

import (
	"flag"
)

type Flags struct {
	RunAddr     string
	RunBaseAddr string
	Parsed      bool
}

func ParseFlags() Flags {
	f := Flags{}
	fs := flag.NewFlagSet("flags", flag.ExitOnError)
	fs.StringVar(&f.RunAddr, "a", "localhost:8080", "address and port to run server, example: localhost:8080")
	fs.StringVar(&f.RunBaseAddr, "b", "http://localhost:8080", "server address")
	fsParsed := fs.Parsed()
	if !fsParsed {
		_ = fs.Parse([]string{})
		f.Parsed = true
	}
	return f
}
