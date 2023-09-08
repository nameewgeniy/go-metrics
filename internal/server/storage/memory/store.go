package memory

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (i MemStorage) AddGauge(name string, value float64) {
	i.gauge[name] = value
}

func (i MemStorage) AddCount(name string, value int64) {
	i.counter[name] += value
}
