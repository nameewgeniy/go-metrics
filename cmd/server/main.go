package main

import (
	"errors"
	"github.com/nameewgeniy/go-metrics/internal/logger"
	"github.com/nameewgeniy/go-metrics/internal/server"
	"github.com/nameewgeniy/go-metrics/internal/server/conf"
	"github.com/nameewgeniy/go-metrics/internal/server/handlers"
	"github.com/nameewgeniy/go-metrics/internal/server/storage/memory"
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

	if err = logger.Initialize(f.logLevel); err != nil {
		return err
	}

	store := memory.NewMemory()
	handler := handlers.NewMuxHandlers(store)

	cnf := conf.NewServerConf(f.addr)
	srv := server.NewServer(cnf, handler)

	sig := make(chan os.Signal, 1)
	defer close(sig)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)

	errorCh := make(chan error)
	defer close(errorCh)

	go func() { errorCh <- srv.Listen() }()

	select {
	case <-sig:
		return errors.New("stop app")
	case err = <-errorCh:
		return err
	}
}
