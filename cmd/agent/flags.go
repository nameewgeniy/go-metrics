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
	pollIntervalSec   int
	reportIntervalSec int
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

	flag.Parse()

	pAddr := os.Getenv("ADDRESS")
	if pAddr != "" {
		f.pushAddress = pAddr
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
