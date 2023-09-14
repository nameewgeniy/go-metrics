package storage

import (
	"errors"
)

var ErrItemNotFound = errors.New("item not found")

type MetricsItemGauage struct {
	Name  string
	Value float64
}

type MetricsItemCounter struct {
	Name  string
	Value int64
}

type Storage interface {
	AddGauage(MetricsItemGauage) error
	FindGauageItem(name string) (MetricsItemGauage, error)
	FindGauageAll() ([]MetricsItemGauage, error)

	AddCounter(MetricsItemCounter) error
	FindCounterItem(name string) (MetricsItemCounter, error)
	FindCounterAll() ([]MetricsItemCounter, error)
}

//
//func (m *MetricsItem) ValueToFloat64() (float64, error) {
//	switch v := m.Value.(type) {
//	case float64:
//		return v, nil
//	case int64:
//		return float64(v), nil
//	case string:
//		val, err := strconv.ParseFloat(v, 64)
//		if err != nil {
//			return 0, err
//		}
//		return val, nil
//	default:
//		return 0, fmt.Errorf("unsupported value type")
//	}
//}
//
//func (m *MetricsItem) ValueToInt64() (int64, error) {
//	switch v := m.Value.(type) {
//	case int64:
//		return v, nil
//	case float64:
//		return int64(v), nil
//	case string:
//		val, err := strconv.ParseInt(v, 10, 64)
//		if err != nil {
//			return 0, err
//		}
//		return val, nil
//	default:
//		return 0, fmt.Errorf("unsupported value type")
//	}
//}
//
//func (m *MetricsItem) String() string {
//	return fmt.Sprintf("%v", m.Value)
//}
