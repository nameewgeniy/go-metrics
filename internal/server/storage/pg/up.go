package pg

import (
	"context"
)

func (p Pg) Up(_ context.Context) error {

	return p.migrationUp()
}
