package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go-metrics/internal/server"
	"go-metrics/internal/server/conf"
	"go-metrics/internal/server/handlers"
	"go-metrics/internal/server/storage/memory"
	"go-metrics/internal/server/storage/pg"
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

	db, err := initDbConnect(f.databaseDsn)
	if err != nil {
		return err
	}

	pgconf := conf.NewPgStorageConf(db)
	pgstore := pg.NewPgStorage(pgconf)

	memcfg := conf.NewMemoryStorageConf(f.fileStoragePath)
	memstore := memory.NewMemoryStorage(memcfg)

	handler := handlers.NewMuxHandlers(memstore, pgstore)

	cnf := conf.NewServerConf(f.addr, f.storeInterval, f.restore)
	srv := server.NewServer(cnf, handler, memstore, memstore)

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

func initDbConnect(databaseDsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func handlePanic(errorCh chan<- error, stop context.CancelFunc) {
	if r := recover(); r != nil {
		errorCh <- fmt.Errorf("panic: %v", r)
		stop()
	}
}
