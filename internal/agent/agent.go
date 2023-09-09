package agent

import (
	"sync"
	"time"
)

type Conf interface {
	PollInterval() time.Duration
	ReportInterval() time.Duration
}

type Metrics interface {
	Sync()
	Push()
}

type Agent struct {
	cnf Conf
	m   Metrics
}

func NewAgent(c Conf, m Metrics) *Agent {
	return &Agent{
		cnf: c,
		m:   m,
	}
}

func (s Agent) Do() error {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.callByTime(s.cnf.PollInterval(), s.m.Sync)
	}()

	go func() {
		defer wg.Done()
		s.callByTime(s.cnf.ReportInterval(), s.m.Push)
	}()

	wg.Wait()

	return nil
}

func (s Agent) callByTime(d time.Duration, f func()) {
	for {
		f()
		time.Sleep(d)
	}
}
