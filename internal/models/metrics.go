package models

import (
	"errors"
	"fmt"
	"go-metrics/internal"
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
	case internal.CounterType:
		return strconv.FormatInt(*m.Delta, 10), nil
	case internal.GaugeType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported metrics type: %s", m.MType)
	}
}

type MetricsFactory struct{}

func NewMetricsFactory() *MetricsFactory {
	return &MetricsFactory{}
}

func (f *MetricsFactory) MakeFromMapForUpdateMetrics(vars map[string]string) (*Metrics, error) {

	if vars["name"] == "" || vars["type"] == "" || vars["value"] == "" {
		return nil, ErrRequiredFields
	}

	m := Metrics{
		ID:    vars["name"],
		MType: vars["type"],
	}

	switch m.MType {
	case internal.CounterType:
		numValue, err := strconv.ParseInt(vars["value"], 10, 64)
		if err != nil {
			return nil, err
		}
		m.Delta = &numValue

	case internal.GaugeType:
		numValue, err := strconv.ParseFloat(vars["value"], 64)
		if err != nil {
			return nil, err
		}
		m.Value = &numValue
	}

	return &m, nil
}

func (f *MetricsFactory) MakeFromMapForGetMetrics(vars map[string]string) (*Metrics, error) {

	if vars["name"] == "" || vars["type"] == "" {
		return nil, ErrRequiredFields
	}

	m := Metrics{
		ID:    vars["name"],
		MType: vars["type"],
	}

	return &m, nil
}

func (f *MetricsFactory) MakeFromBytesForUpdateMetrics(bytes []byte) (*Metrics, error) {

	m := Metrics{}

	if err := m.UnmarshalJSON(bytes); err != nil {
		return nil, err
	}

	if m.ID == "" || m.MType == "" || (nil == m.Value && nil == m.Delta) {
		return nil, ErrRequiredFields
	}

	if nil != m.Value && nil != m.Delta {
		return nil, ErrMutuallyExclusiveValues
	}

	return &m, nil
}

func (f *MetricsFactory) MakeFromBytesForGetMetrics(bytes []byte) (*Metrics, error) {

	m := Metrics{}

	if err := m.UnmarshalJSON(bytes); err != nil {
		return nil, err
	}

	if m.ID == "" || m.MType == "" {
		return nil, ErrRequiredFields
	}

	return &m, nil
}
