package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
)

type flags struct {
	addr     string
	logLevel string
}

func (f *flags) validate() error {
	if _, _, err := net.SplitHostPort(f.addr); err != nil {
		return fmt.Errorf("address is not valid: %s", err)
	}

	if _, err := zap.ParseAtomicLevel(f.logLevel); err != nil {
		return fmt.Errorf("log level is not valid: %s", err)
	}

	return nil
}

func parseFlags() (*flags, error) {
	var f flags

	flag.StringVar(&f.addr, "a", "localhost:8080", "address")
	flag.StringVar(&f.logLevel, "l", "info", "log level")
	flag.Parse()

	envAddr := os.Getenv("ADDRESS")
	if envAddr != "" {
		f.addr = envAddr
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		f.logLevel = envLogLevel
	}

	if err := f.validate(); err != nil {
		return nil, err
	}

	return &f, nil
}
