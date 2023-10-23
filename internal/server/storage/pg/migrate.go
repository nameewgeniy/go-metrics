package pg

import (
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (p Pg) migrationUp() error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(p.mDialect); err != nil {
		return err
	}

	if err := goose.Up(p.c.DB(), p.mDir); err != nil {
		return err
	}

	return nil
}

func (p Pg) migrationDown() error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(p.mDialect); err != nil {
		return err
	}

	if err := goose.Down(p.c.DB(), p.mDir); err != nil {
		return err
	}

	return nil
}
