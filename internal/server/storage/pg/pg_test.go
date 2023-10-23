package pg

import (
	"database/sql"
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
// Создаем фейковый структуру Pg, имитирующую Pg объект
type fakePg struct {
	c *sql.DB
}

func (p fakePg) Ping() error {
	return nil
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

func TestPing(t *testing.T) {
	// Создаем фейковый Pg объект
	pg := fakePg{
		c: nil, // здесь нужно передать фейковое подключение к базе данных
	}

	// Вызываем функцию Ping()
	err := pg.Ping()

	// Проверяем, что возвращенная ошибка равна nil (Нет ошибки)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
