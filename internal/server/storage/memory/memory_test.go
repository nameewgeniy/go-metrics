package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory_Add(t *testing.T) {
	m := NewMemory()

	// Тестирование добавления счетчика
	counter := storage.MetricsItemCounter{
		Name:  "requests",
		Value: int64(10),
	}
	err := m.AddCounter(counter)
	assert.NoError(t, err)
	v, _ := m.Counter.Load("requests")
	assert.Equal(t, int64(10), v.(int64))

	err = m.AddCounter(counter)
	assert.NoError(t, err)
	v, _ = m.Counter.Load("requests")
	assert.Equal(t, int64(20), v.(int64))

	// Тестирование добавления метрики
	gauge := storage.MetricsItemGauge{
		Name:  "temperature",
		Value: 25.5,
	}
	err = m.AddGauge(gauge)
	assert.NoError(t, err)
	v, _ = m.Gauge.Load("temperature")
	assert.Equal(t, 25.5, v.(float64))
}

func TestMemory_FindCounterItem(t *testing.T) {
	m := &Memory{}
	m.Counter.Store("counter1", int64(10))

	// Проверяем, что возвращено правильное значение для существующего элемента
	item, err := m.FindCounterItem("counter1")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), item.Value) // ожидаемое значение

	// Проверяем, что возвращается пустой элемент для несуществующего имени
	item, err = m.FindCounterItem("counter2")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), item.Value) // ожидаемое значение
}

func TestMemory_FindGaugeItem(t *testing.T) {
	m := &Memory{}
	m.Gauge.Store("gauge1", float64(10))

	// Проверяем, что возвращено правильное значение для существующего элемента
	item, err := m.FindGaugeItem("gauge1")
	assert.NoError(t, err)
	assert.Equal(t, float64(10), item.Value) // ожидаемое значение

	// Проверяем, что возвращается пустой элемент для несуществующего имени
	item, err = m.FindGaugeItem("gauge2")
	assert.NoError(t, err)
	assert.Equal(t, float64(0), item.Value) // ожидаемое значение
}
