package pg

import (
	"context"
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
	"time"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

func (p Pg) Migrate() error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(p.c.Db(), "migrations"); err != nil {
		return err
	}
}
