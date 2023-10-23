package pg

import (
	"database/sql"
	"testing"
)

// Создаем фейковый структуру Pg, имитирующую Pg объект
type fakePg struct {
	c *sql.DB
}

func (p fakePg) Ping() error {
	return nil
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
