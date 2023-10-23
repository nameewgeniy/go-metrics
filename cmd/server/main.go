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

	if err = logger.Initialize(f.logLevel); err != nil {
		return err
	}

	var strg storage.Storage
	var snapshot memory.SnapshotStorage
	var ping handlers.Ping

	memcfg := conf.NewMemoryStorageConf(f.fileStoragePath, f.restore)
	mstrg := memory.NewMemoryStorage(memcfg)
	strg, snapshot = mstrg, mstrg

	if f.databaseDsn != "" {

		pgconn, err := sql.Open("pgx", f.databaseDsn)
		defer func() {
			_ = pgconn.Close()
		}()

		if err != nil {
			return err
		}

		pgconf := conf.NewPgStorageConf(pgconn, f.downMigrations)
		pstrg := pg.NewPgStorage(pgconf)
		strg, ping = pstrg, pstrg
	}

	handler := handlers.NewMuxHandlers(strg, ping)

	cnf := conf.NewServerConf(f.addr, f.storeInterval, f.restore)
	srv := server.NewServer(cnf, handler, strg, snapshot)

	return srv.Run()
}
