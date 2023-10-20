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

	_, err := p.FindGaugeItem(gauge.Name)

	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			return p.insertGauge(ctx, gauge)
		}

		return err
	}

	return p.updateGauge(ctx, gauge)
}

func (p Pg) insertGauge(ctx context.Context, gauge storage.MetricsItemGauge) error {
	baseQuery := "INSERT INTO #table# (name, value) VALUES ($1, $2)"
	preparedQuery := strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)

	_, err := p.c.Db().ExecContext(ctx, preparedQuery, gauge.Name, gauge.Value)

	return err
}

func (p Pg) updateGauge(ctx context.Context, gauge storage.MetricsItemGauge) error {
	baseQuery := "UPDATE #table# SET value = $2 WHERE name = $1"
	preparedQuery := strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)

	_, err := p.c.Db().ExecContext(ctx, preparedQuery, gauge.Name, gauge.Value)

	return err
}

func (p Pg) FindGaugeItem(name string) (storage.MetricsItemGauge, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res := storage.MetricsItemGauge{}
	baseQuery := "SELECT name, value FROM #table# WHERE name = $1"
	preparedQuery := strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)

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

func (p Pg) FindGaugeAll() ([]storage.MetricsItemGauge, error) {

	var res []storage.MetricsItemGauge

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	baseQuery := "SELECT name, value FROM #table#"
	preparedQuery := strings.NewReplacer("#table#", p.gaugeTableName).Replace(baseQuery)

	rows, err := p.c.Db().QueryContext(ctx, preparedQuery)
	if err != nil {
		return res, err
	}
	defer rows.Close()

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
