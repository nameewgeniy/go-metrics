package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

type flags struct {
	pushAddress       string
	hashKey           string
	pollIntervalSec   int
	reportIntervalSec int
	rateLimit         int
}

func (f *flags) validate() error {
	if _, _, err := net.SplitHostPort(f.pushAddress); err != nil {
		return fmt.Errorf("address is not valid: %s", err)
	}

	if f.pollIntervalSec < 1 {
		return fmt.Errorf("poll interval must be greater than 1, current: %d", f.pollIntervalSec)
	}

	if f.reportIntervalSec < 1 {
		return fmt.Errorf("report interval must be greater than 1, current: %d", f.reportIntervalSec)
	}

	return nil
}

func parseFlags() (*flags, error) {
	var f flags

	flag.StringVar(&f.pushAddress, "a", "localhost:8080", "push address")
	flag.IntVar(&f.pollIntervalSec, "p", 2, "poll interval in seconds")
	flag.IntVar(&f.reportIntervalSec, "r", 10, "report interval in seconds")
	flag.IntVar(&f.rateLimit, "l", 1, "rate limit")
	flag.StringVar(&f.hashKey, "k", "", "hash key")

	flag.Parse()

	pAddr := os.Getenv("ADDRESS")
	if pAddr != "" {
		f.pushAddress = pAddr
	}

	k := os.Getenv("KEY")
	if k != "" {
		f.hashKey = k
	}

	rl := os.Getenv("RATE_LIMIT")
	if rl != "" {
		v, err := strconv.Atoi(rl)
		if err != nil {
			return nil, fmt.Errorf("RATE_LIMIT is not valid: %s", err)
		}
		f.rateLimit = v
	}

	pi := os.Getenv("POLL_INTERVAL")
	if pi != "" {
		v, err := strconv.Atoi(pi)
		if err != nil {
			return nil, fmt.Errorf("POLL_INTERVAL is not valid: %s", err)
		}
		f.pollIntervalSec = v
	}

	ri := os.Getenv("REPORT_INTERVAL")
	if ri != "" {
		v, err := strconv.Atoi(ri)
		if err != nil {
			return nil, fmt.Errorf("REPORT_INTERVAL is not valid: %s", err)
		}
		f.reportIntervalSec = v
	}

	if err := f.validate(); err != nil {
		return nil, err
	}

	return &f, nil
}
