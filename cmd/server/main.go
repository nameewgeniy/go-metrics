package main

import (
	"errors"
	"go-metrics/internal/logger"
	"go-metrics/internal/server"
	"go-metrics/internal/server/conf"
	"go-metrics/internal/server/handlers"
	"go-metrics/internal/server/storage/memory"
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

	mcfg := conf.NewStorageConf(f.fileStoragePath)
	store := memory.NewMemory(mcfg)

	handler := handlers.NewMuxHandlers(store)

	cnf := conf.NewServerConf(f.addr, f.storeInterval, f.restore)
	srv := server.NewServer(cnf, handler, store)

	sig := make(chan os.Signal, 1)
	defer close(sig)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)

	errorCh := make(chan error)
	defer close(errorCh)

	go func() { srv.Workers(errorCh, sig) }()
	go func() { srv.Listen(errorCh) }()

	select {
	case <-sig:
		return errors.New("stop app")
	case err = <-errorCh:
		return err
	}
}
