package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	memory := Memory{}

	t.Run("Save gauge metrics", func(t *testing.T) {
		item := storage.MetricsItem{
			Gauge: map[string]float64{
				"metric1": 2.5,
				"metric2": 3.7,
			},
			Counter: map[string]int64{},
		}

		err := memory.Save(item)

		assert.Nil(t, err)

		// Check if gauge metrics were saved correctly
		expectedGauges := map[string]float64{
			"metric1": 2.5,
			"metric2": 3.7,
		}

		for name, value := range expectedGauges {
			actual, ok := memory.gauge.Load(name)
			assert.True(t, ok, "Expected gauge metric %s to be saved", name)
			assert.Equal(t, value, actual, "Expected gauge metric %s to have value %f", name, value)
		}

	})
}
