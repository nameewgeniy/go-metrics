package pg

import "context"

func (p Pg) Down(_ context.Context) error {

	if p.c.DownMigrations() {
		return p.migrationDown()
	}

	return nil
}
