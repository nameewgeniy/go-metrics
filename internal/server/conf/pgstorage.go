package conf

import "database/sql"

type PgStorageConf struct {
	db *sql.DB
}

func NewPgStorageConf(db *sql.DB) *PgStorageConf {
	return &PgStorageConf{
		db: db,
	}
}

func (c PgStorageConf) DB() *sql.DB {
	return c.db
}
