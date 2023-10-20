package pg

import (
	"context"
	"database/sql"
	"errors"
	"go-metrics/internal/server/storage"
	"strings"
	"time"
)

func (p Pg) AddCounter(counter storage.MetricsItemCounter) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	existCounter, err := p.FindCounterItem(counter.Name)

	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			return p.insertCounter(ctx, counter)
		}

		return err
	}

	counter.Value += existCounter.Value
	return p.updateCounter(ctx, counter)
}

func (p Pg) insertCounter(ctx context.Context, counter storage.MetricsItemCounter) error {
	baseQuery := "INSERT INTO #table# (name, value) VALUES ($1, $2)"
	preparedQuery := strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)

	_, err := p.c.Db().ExecContext(ctx, preparedQuery, counter.Name, counter.Value)

	return err
}

func (p Pg) updateCounter(ctx context.Context, counter storage.MetricsItemCounter) error {
	baseQuery := "UPDATE #table# SET value = $2 WHERE name = $1"
	preparedQuery := strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)

	_, err := p.c.Db().ExecContext(ctx, preparedQuery, counter.Name, counter.Value)

	return err
}

func (p Pg) FindCounterItem(name string) (storage.MetricsItemCounter, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res := storage.MetricsItemCounter{}
	baseQuery := "SELECT name, value FROM #table# WHERE name = $1"
	preparedQuery := strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)

	rows := p.c.Db().QueryRowContext(ctx, preparedQuery, name)

	if err := rows.Scan(&res.Name, &res.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, storage.ErrItemNotFound
		}

		return res, err
	}

	if err := rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (p Pg) FindCounterAll() ([]storage.MetricsItemCounter, error) {
	var res []storage.MetricsItemCounter

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	baseQuery := "SELECT name, value FROM #table#"
	preparedQuery := strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)

	rows, err := p.c.Db().QueryContext(ctx, preparedQuery)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var item storage.MetricsItemCounter
		if err := rows.Scan(&item.Name, &item.Value); err != nil {
			return res, err
		}
		res = append(res, item)
	}

	if err := rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
