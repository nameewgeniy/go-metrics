package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

type flags struct {
	addr string
}

func (f *flags) validate() error {
	if _, _, err := net.SplitHostPort(f.addr); err != nil {
		return fmt.Errorf("address is not valid: %s", err)
	}

	return nil
}

func parseFlags() (*flags, error) {
	var f flags

	flag.StringVar(&f.addr, "a", "localhost:8080", "address")
	flag.Parse()

	addr := os.Getenv("ADDRESS")
	if addr != "" {
		f.addr = addr
	}

	if err := f.validate(); err != nil {
		return nil, err
	}

	return &f, nil
}
