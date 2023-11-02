package metrics

import (
	"encoding/json"
	"go-metrics/internal/shared"
	"strconv"
)

type MetricsFactory struct{}

func NewMetricsFactory() *MetricsFactory {
	return &MetricsFactory{}
}

func (f *MetricsFactory) MakeFromMapForUpdateMetrics(vars map[string]string) (*Metrics, error) {

	m := Metrics{
		ID:    vars["name"],
		MType: vars["type"],
	}

	switch m.MType {
	case shared.CounterType:
		numValue, err := strconv.ParseInt(vars["value"], 10, 64)
		if err != nil {
			return nil, err
		}
		m.Delta = &numValue

	case shared.GaugeType:
		numValue, err := strconv.ParseFloat(vars["value"], 64)
		if err != nil {
			return nil, err
		}
		m.Value = &numValue
	}

	if err := m.ValidateState(); err != nil {
		return nil, err
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

	if err := m.ValidateState(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (f *MetricsFactory) MakeFromBytesForBatchUpdateMetrics(bytes []byte) ([]Metrics, error) {

	var m []Metrics

	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	for i := range m {
		if m[i].ID == "" || m[i].MType == "" || (nil == m[i].Value && nil == m[i].Delta) {
			return nil, ErrRequiredFields
		}

		if nil != m[i].Value && nil != m[i].Delta {
			return nil, ErrMutuallyExclusiveValues
		}
	}

	return m, nil
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
