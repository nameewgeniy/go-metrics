package service

import (
	"fmt"
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"strconv"
)

const gaugeType = "gauge"
const counterType = "counter"

type Storage interface {
	Save(storage.MetricsItem) error
}

type Metrics struct {
	s Storage
}

func NewMetrics(s Storage) *Metrics {
	return &Metrics{
		s: s,
	}
}

func (m Metrics) Update(mType, mName, mValue string) error {

	i := storage.MetricsItem{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}

	switch mType {
	case gaugeType:

		mValueNumber, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return err
		}
		i.Gauge[mName] = mValueNumber

	case counterType:

		mValueNumber, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			return err
		}
		i.Counter[mName] = mValueNumber

	default:
		return fmt.Errorf("type not supported")
	}

	err := m.s.Save(i)

	if err != nil {
		return err
	}

	return nil
}
