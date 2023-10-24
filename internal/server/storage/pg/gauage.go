package pg

import (
	"context"
	"database/sql"
	"errors"
	"go-metrics/internal/server/storage"
	"strings"
	"time"
)

func (p Pg) AddGauge(gauge storage.MetricsItemGauge) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.c.DB().ExecContext(ctx, p.upsertGaugeSQL(), gauge.Name, gauge.Value)

	return err
}

func (p Pg) AddBatchGauges(gauges []storage.MetricsItemGauge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tr, err := p.c.DB().BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	for _, gauge := range gauges {
		if _, err = tr.ExecContext(ctx, p.upsertGaugeSQL(), gauge.Name, gauge.Value); err != nil {
			return err
		}
	}

	return tr.Commit()
}

func (p Pg) upsertGaugeSQL() string {
	baseQuery := "INSERT INTO #table# (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2"
	return strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)
}

func (p Pg) FindGaugeItem(name string) (storage.MetricsItemGauge, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res := storage.MetricsItemGauge{}
	baseQuery := "SELECT name, value FROM #table# WHERE name = $1"
	preparedQuery := strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)

	rows := p.c.DB().QueryRowContext(ctx, preparedQuery, name)

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

func (p Pg) FindGaugeAll() ([]storage.MetricsItemGauge, error) {

	var res []storage.MetricsItemGauge

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	baseQuery := "SELECT name, value FROM #table#"
	preparedQuery := strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)

	rows, err := p.c.DB().QueryContext(ctx, preparedQuery)
	if err != nil {
		return res, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var item storage.MetricsItemGauge
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
