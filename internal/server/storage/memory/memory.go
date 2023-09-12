package memory

import (
	"fmt"
	"github.com/nameewgeniy/go-metrics/internal"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"sync"
)

type Memory struct {
	Gauge   sync.Map
	Counter sync.Map
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Add(i storage.MetricsItem) error {

	switch i.Type {
	case internal.CounterType:
		v, err := i.ValueToInt64()

		if err != nil {
			return err
		}

		m.addCount(i.Name, v)
	case internal.GaugeType:
		v, err := i.ValueToFloat64()

		if err != nil {
			return err
		}

		m.addGauge(i.Name, v)
	default:
		return fmt.Errorf("unsupported type")
	}

	return nil
}

func (m *Memory) addGauge(name string, value float64) {
	m.Gauge.Store(name, value)
}

func (m *Memory) addCount(name string, value int64) {
	if oldValue, ok := m.Counter.Load(name); ok {
		value += oldValue.(int64)
	}

	m.Counter.Store(name, value)
}

func (m *Memory) Find(mType, name string) (storage.MetricsItem, error) {

	res := storage.MetricsItem{
		Type: mType,
		Name: name,
	}

	switch mType {
	case internal.CounterType:
		if val, ok := m.Counter.Load(name); ok {
			res.Value = val
			return res, nil
		}
	case internal.GaugeType:
		if val, ok := m.Gauge.Load(name); ok {
			res.Value = val
			return res, nil
		}
	default:
		return res, fmt.Errorf("unsupported type")
	}

	return res, storage.ErrItemNotFound
}

func (m *Memory) All() ([]storage.MetricsItem, error) {
	var res []storage.MetricsItem

	m.Counter.Range(func(name, value interface{}) bool {
		res = append(res, storage.MetricsItem{
			Type:  internal.CounterType,
			Name:  name.(string),
			Value: value,
		})
		return true
	})

	m.Gauge.Range(func(name, value interface{}) bool {
		res = append(res, storage.MetricsItem{
			Type:  internal.GaugeType,
			Name:  name.(string),
			Value: value,
		})
		return true
	})

	return res, nil
}
