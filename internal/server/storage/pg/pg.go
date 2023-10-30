package pg

import (
	"database/sql"
)

type PgStorageConfig interface {
	DB() *sql.DB
	DownMigrations() bool
}

type Pg struct {
	c                PgStorageConfig
	mDialect         string
	mDir             string
	gaugeTableName   string
	counterTableName string
}

func NewPgStorage(c PgStorageConfig) *Pg {
	return &Pg{
		c:                c,
		mDialect:         "postgres",
		mDir:             "migrations",
		gaugeTableName:   "metrics_gauge",
		counterTableName: "metrics_counter",
	}
}
