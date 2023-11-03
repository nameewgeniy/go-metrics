package metrics

import (
	"errors"
	"fmt"
	"go-metrics/internal/shared"
	"strconv"
)

var ErrRequiredFields = errors.New("one or more required fields are empty")
var ErrMutuallyExclusiveValues = errors.New("mutually exclusive values are passed")

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) ValueByType() (string, error) {
	switch m.MType {
	case shared.CounterType:
		return strconv.FormatInt(*m.Delta, 10), nil
	case shared.GaugeType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported metrics type: %s", m.MType)
	}
}

func (m Metrics) ValidateState() error {

	if m.ID == "" || m.MType == "" {
		return ErrRequiredFields
	}

	if m.MType == shared.GaugeType && nil == m.Value {
		return fmt.Errorf("validate: field value is empty:  %w", ErrRequiredFields)
	}

	if m.MType == shared.CounterType && nil == m.Delta {
		return fmt.Errorf("validate: field delta is empty:  %w", ErrRequiredFields)
	}

	if nil != m.Value && nil != m.Delta {
		return ErrMutuallyExclusiveValues
	}

	return nil
}
