package conf

import "database/sql"

type PgStorageConf struct {
	db             *sql.DB
	downMigrations bool
}

func NewPgStorageConf(db *sql.DB, downMigrations bool) *PgStorageConf {
	return &PgStorageConf{
		db:             db,
		downMigrations: downMigrations,
	}
}

func (c PgStorageConf) Db() *sql.DB {
	return c.db
}

func (c PgStorageConf) DownMigrations() bool {
	return c.downMigrations
}
