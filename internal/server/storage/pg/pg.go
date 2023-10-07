package pg

import (
	"context"
	"database/sql"
	"time"
)

type PgStorageConfig interface {
	Db() *sql.DB
}

type Pg struct {
	c PgStorageConfig
}

func NewPgStorage(c PgStorageConfig) *Pg {
	return &Pg{
		c: c,
	}
}

func (p Pg) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return p.c.Db().PingContext(ctx)
}
