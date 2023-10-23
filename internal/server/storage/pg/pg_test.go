package pg

import (
	"testing"
)

func TestUpsertCounterSql(t *testing.T) {
	p := Pg{
		counterTableName: "counters",
	}
	expectedQuery := "INSERT INTO counters as t (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = t.value + excluded.value"
	actualQuery := p.upsertCounterSql()
	if actualQuery != expectedQuery {
		t.Errorf("Expected query: %s, but got: %s", expectedQuery, actualQuery)
	}
}

func TestReturnsCorrectSQLQuery(t *testing.T) {
	p := Pg{
		gaugeTableName: "gauge_table",
	}
	expectedQuery := "INSERT INTO gauge_table (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2"
	actualQuery := p.upsertGaugeSql()
	if actualQuery != expectedQuery {
		t.Errorf("Expected query: %s, but got: %s", expectedQuery, actualQuery)
	}
}
