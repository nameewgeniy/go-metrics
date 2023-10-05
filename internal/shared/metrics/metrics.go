package metrics

import (
	"fmt"
	"go-metrics/internal"
	"strconv"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) ValueByType() (string, error) {
	switch m.MType {
	case internal.CounterType:
		return strconv.FormatInt(*m.Delta, 10), nil
	case internal.GaugeType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported metrics type: %s", m.MType)
	}
}
