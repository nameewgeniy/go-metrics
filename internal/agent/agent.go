package agent

import (
	"context"
	"fmt"
	"time"
)

type Metrics interface {
	Sync()
	Push()
}

type Config interface {
	PollInterval() time.Duration
	ReportInterval() time.Duration
}

type Agent struct {
	cnf Config
	m   Metrics
}

func NewAgent(c Config, m Metrics) *Agent {
	return &Agent{
		cnf: c,
		m:   m,
	}
}

func (s Agent) Do(ctx context.Context, errorCh chan<- error) {

	pollTicker := time.NewTicker(s.cnf.PollInterval())
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(s.cnf.ReportInterval())
	defer reportTicker.Stop()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errorCh <- fmt.Errorf("report panic: %v", r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				s.m.Sync()
			}
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errorCh <- fmt.Errorf("sync panic: %v", r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-reportTicker.C:
				s.m.Push()
			}
		}
	}()
}
