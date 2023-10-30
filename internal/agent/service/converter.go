package service

import (
	"math"
)

// MetricsValueUintToInt кидает панику в случае, если значение в метрике uint64 больше int64
// кидается паника, так как это исключительная ситуация
func MetricsValueUintToInt(value uint64) int64 {
	if value > uint64(math.MaxInt64) {
		panic("the value exceeds the maximum allowed value int64")
	}

	return int64(value)
}
