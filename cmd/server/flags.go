package main

import (
	"flag"
)

type flags struct {
	addr string
}

func parseFlags() (*flags, error) {
	var f flags

	flag.StringVar(&f.addr, "a", ":8080", "address")
	flag.Parse()

	return &f, nil
}
