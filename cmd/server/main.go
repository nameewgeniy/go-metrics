package main

import (
	"context"
	"fmt"
	"go-metrics/internal/server"
	"go-metrics/internal/server/conf"
	"go-metrics/internal/server/handlers"
	"go-metrics/internal/server/storage/memory"
	"go-metrics/internal/shared/logger"
	"golang.org/x/sync/errgroup"
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
	srv := server.NewServer(cnf, handler, store, store)

	err = srv.Restore()
	if err != nil {
		return err
	}

	errorCh := make(chan error)
	defer close(errorCh)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eg, errCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer handlePanic(errorCh, cancel)
		return srv.Snapshot(errCtx)
	})

	eg.Go(func() error {
		defer handlePanic(errorCh, cancel)
		return srv.Listen(errCtx)
	})

	go func() {
		errorCh <- eg.Wait()
	}()

	return <-errorCh
}

func handlePanic(errorCh chan<- error, stop context.CancelFunc) {
	if r := recover(); r != nil {
		errorCh <- fmt.Errorf("panic: %v", r)
		stop()
	}
}
