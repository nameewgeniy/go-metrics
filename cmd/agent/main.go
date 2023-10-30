package main

import (
	"context"
	"errors"
	"go-metrics/internal/agent"
	"go-metrics/internal/agent/conf"
	"go-metrics/internal/agent/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	f, err := parseFlags()
	if err != nil {
		return err
	}

	scf := conf.NewSenderConfig(f.pushAddress)
	snd := service.NewMetricSender(scf)
	rm := service.NewRuntimeMetrics(snd)

	cf := conf.NewAgentConf(f.pollIntervalSec, f.reportIntervalSec)
	a := agent.NewAgent(cf, rm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errorCh := make(chan error)
	defer close(errorCh)

	sig := make(chan os.Signal, 1)
	defer close(sig)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)

	a.Do(ctx, errorCh)

	select {
	case <-sig:
		return errors.New("stop app")
	case err = <-errorCh:
		return err
	}
}
