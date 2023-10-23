package pg

import (
	"context"
	"time"
)

func (p Pg) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return p.c.DB().PingContext(ctx)
}
