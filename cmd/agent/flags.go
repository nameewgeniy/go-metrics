package main

import (
	"flag"
)

type flags struct {
	pushAddress       string
	pollIntervalSec   int
	reportIntervalSec int
}

func parseFlags() (*flags, error) {
	var f flags

	flag.StringVar(&f.pushAddress, "a", "localhost:8080", "push address")
	flag.IntVar(&f.pollIntervalSec, "p", 2, "poll interval in seconds")
	flag.IntVar(&f.reportIntervalSec, "r", 10, "report interval in seconds")

	flag.Parse()

	return &f, nil
}
