package agent

import (
	"context"
	"fmt"
	"time"
)

type Metrics interface {
	Sync()
	Push(ctx context.Context)
}

type Config interface {
	PollInterval() time.Duration
	ReportInterval() time.Duration
	RateLimit() int
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

func (a Agent) Do(ctx context.Context, errorCh chan<- error) {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errorCh <- fmt.Errorf("report panic: %v", r)
			}
		}()

		pollTicker := time.NewTicker(a.cnf.PollInterval())
		defer pollTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				a.m.Sync()
			}
		}
	}()

	for i := 0; i < a.cnf.RateLimit(); i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					errorCh <- fmt.Errorf("sync panic: %v", r)
				}
			}()

			ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			reportTicker := time.NewTicker(a.cnf.ReportInterval())
			defer reportTicker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-reportTicker.C:
					a.m.Push(ctxTimeout)
				}
			}
		}()
	}
}
