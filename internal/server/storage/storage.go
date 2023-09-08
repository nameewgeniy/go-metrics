package storage

type MetricsItem struct {
	Gauge   map[string]float64
	Counter map[string]int64
}
