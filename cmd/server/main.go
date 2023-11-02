package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go-metrics/internal/server"
	"go-metrics/internal/server/conf"
	"go-metrics/internal/server/handlers"
	"go-metrics/internal/server/storage"
	"go-metrics/internal/server/storage/memory"
	"go-metrics/internal/server/storage/pg"
	"go-metrics/internal/shared/logger"
	"go-metrics/internal/shared/signature"
	"log"
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

	if err = logger.Singleton(f.logLevel); err != nil {
		return err
	}

	signature.Singleton(f.hashKey)

	var conn *sql.DB
	if f.databaseDsn != "" {

		if conn, err = sql.Open("pgx", f.databaseDsn); err != nil {
			return err
		}

		defer func() {
			_ = conn.Close()
		}()
	}

	mainStorage, snapshot, ping, err := newStorageServices(f, conn)
	if err != nil {
		return err
	}

	handler := handlers.NewMuxHandlers(mainStorage, ping)

	cnf := conf.NewServerConf(f.addr, f.storeInterval, f.restore)
	srv := server.NewServer(cnf, handler, mainStorage, snapshot)

	return srv.Run()
}

func newStorageServices(f *flags, conn *sql.DB) (storage.Storage, memory.SnapshotStorage, handlers.Ping, error) {
	var store storage.Storage
	var snapshot memory.SnapshotStorage
	var ping handlers.Ping

	// init memory storage
	memStorage := memory.NewMemoryStorage(
		conf.NewMemoryStorageConf(f.fileStoragePath, f.restore),
	)

	// memStorage реализует интерфейс Storage и SnapshotStorage
	store, snapshot = memStorage, memStorage

	if conn != nil {
		pgStorage := pg.NewPgStorage(
			conf.NewPgStorageConf(conn, f.downMigrations),
		)

		// pgStorage реализует интерфейс Storage и Ping
		store, ping = pgStorage, pgStorage
	}

	return store, snapshot, ping, nil
}
