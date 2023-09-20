package service

import "fmt"

func MetricsValueToString[T float64 | uint32 | uint64](value T) string {
	return fmt.Sprintf("%v", value)
}
