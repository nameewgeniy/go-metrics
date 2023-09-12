package memory

import (
	"github.com/nameewgeniy/go-metrics/internal"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory_Add(t *testing.T) {
	m := NewMemory()

	// Тестирование добавления счетчика
	counter := storage.MetricsItem{
		Name:  "requests",
		Type:  internal.CounterType,
		Value: int64(10),
	}
	err := m.Add(counter)
	assert.NoError(t, err)
	v, _ := m.Counter.Load("requests")
	assert.Equal(t, int64(10), v.(int64))

	// Тестирование добавления метрики
	gauge := storage.MetricsItem{
		Name:  "temperature",
		Type:  internal.GaugeType,
		Value: 25.5,
	}
	err = m.Add(gauge)
	assert.NoError(t, err)
	v, _ = m.Gauge.Load("temperature")
	assert.Equal(t, 25.5, v.(float64))

	// Тестирование неподдерживаемого типа
	unknown := storage.MetricsItem{
		Name:  "unknown",
		Type:  "unknown",
		Value: 42,
	}
	err = m.Add(unknown)
	assert.Error(t, err)
	assert.EqualError(t, err, "unsupported type")
}
