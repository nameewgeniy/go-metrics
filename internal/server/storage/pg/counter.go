package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-metrics/internal/server/storage"
	"strings"
	"time"
)

func (p Pg) AddCounter(counter storage.MetricsItemCounter) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.c.DB().ExecContext(ctx, p.upsertCounterSQL(), counter.Name, counter.Value)

	if err != nil {
		return fmt.Errorf("pg: AddCounter: %w", err)
	}

	return nil
}

func (p Pg) AddBatchCounters(counters []storage.MetricsItemCounter) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tr, err := p.c.DB().BeginTx(ctx, nil)
	defer func() {
		_ = tr.Rollback()
	}()

	if err != nil {
		return fmt.Errorf("pg: AddBatchCounters: Begin Transaction: %w", err)
	}

	for _, counter := range counters {
		if _, err = tr.ExecContext(ctx, p.upsertCounterSQL(), counter.Name, counter.Value); err != nil {
			return fmt.Errorf("pg: AddBatchCounters: %w", err)
		}
	}

	return tr.Commit()
}

func (p Pg) upsertCounterSQL() string {
	baseQuery := "INSERT INTO #table# as t (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = t.value + excluded.value"
	return strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)
}

func (p Pg) FindCounterItem(name string) (storage.MetricsItemCounter, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res := storage.MetricsItemCounter{}
	baseQuery := "SELECT name, value FROM #table# WHERE name = $1"
	preparedQuery := strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)

	rows := p.c.DB().QueryRowContext(ctx, preparedQuery, name)

	if err := rows.Scan(&res.Name, &res.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, storage.ErrItemNotFound
		}

		return res, fmt.Errorf("pg: FindCounterItem: name=%s: %w", name, err)
	}

	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("pg: FindCounterItem: Rows: name=%s: %w", name, err)
	}

	return res, nil
}

func (p Pg) FindCounterAll() ([]storage.MetricsItemCounter, error) {
	var res []storage.MetricsItemCounter

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	baseQuery := "SELECT name, value FROM #table#"
	preparedQuery := strings.NewReplacer("#table#", p.counterTableName).Replace(baseQuery)

	rows, err := p.c.DB().QueryContext(ctx, preparedQuery)
	if err != nil {
		return res, fmt.Errorf("pg: FindCounterAll: Rows: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var item storage.MetricsItemCounter
		if err := rows.Scan(&item.Name, &item.Value); err != nil {
			return res, err
		}
		res = append(res, item)
	}

	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("pg: FindCounterAll: Rows: %w", err)
	}

	return res, nil
}
