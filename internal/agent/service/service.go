package service

import (
	"go-metrics/internal/shared"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type MetricsSender interface {
	SendMemStatsMetric(metricType, metricName, metricValue string) error
}

// Метрки, которые мы будем отправлять из пакета runtime
type metricsTracked struct {
	mType string
	name  string
	value string
}

type extMemStats struct {
	m           *runtime.MemStats
	PollCount   uint64
	RandomValue float64
	mutex       sync.RWMutex
}

type RuntimeMetrics struct {
	memStats *extMemStats
	s        MetricsSender
}

func NewRuntimeMetrics(sender MetricsSender) *RuntimeMetrics {
	return &RuntimeMetrics{
		memStats: &extMemStats{
			m: &runtime.MemStats{},
		},
		s: sender,
	}
}

func (m RuntimeMetrics) Push() {

	m.memStats.mutex.RLock()
	defer m.memStats.mutex.RUnlock()

	for _, v := range m.MetricsTracked() {
		item := v

		go func() {
			err := m.s.SendMemStatsMetric(item.mType, item.name, item.value)
			if err != nil {
				log.Print(err)
			}
		}()
	}
}

func (m RuntimeMetrics) MetricsTracked() []metricsTracked {
	return []metricsTracked{
		{
			mType: shared.GaugeType,
			name:  "Alloc",
			value: MetricsValueToString(m.memStats.m.Alloc),
		},
		{
			mType: shared.GaugeType,
			name:  "NextGC",
			value: MetricsValueToString(m.memStats.m.NextGC),
		},
		{
			mType: shared.GaugeType,
			name:  "BuckHashSys",
			value: MetricsValueToString(m.memStats.m.BuckHashSys),
		},
		{
			mType: shared.GaugeType,
			name:  "Frees",
			value: MetricsValueToString(m.memStats.m.Frees),
		},
		{
			mType: shared.GaugeType,
			name:  "GCCPUFraction",
			value: MetricsValueToString(m.memStats.m.GCCPUFraction),
		},
		{
			mType: shared.GaugeType,
			name:  "Mallocs",
			value: MetricsValueToString(m.memStats.m.Mallocs),
		},
		{
			mType: shared.GaugeType,
			name:  "MSpanSys",
			value: MetricsValueToString(m.memStats.m.MSpanSys),
		},
		{
			mType: shared.GaugeType,
			name:  "MSpanInuse",
			value: MetricsValueToString(m.memStats.m.MSpanInuse),
		},
		{
			mType: shared.GaugeType,
			name:  "MCacheSys",
			value: MetricsValueToString(m.memStats.m.MCacheSys),
		},
		{
			mType: shared.GaugeType,
			name:  "MCacheInuse",
			value: MetricsValueToString(m.memStats.m.MCacheInuse),
		},
		{
			mType: shared.GaugeType,
			name:  "Lookups",
			value: MetricsValueToString(m.memStats.m.Lookups),
		},
		{
			mType: shared.GaugeType,
			name:  "LastGC",
			value: MetricsValueToString(m.memStats.m.LastGC),
		},
		{
			mType: shared.GaugeType,
			name:  "HeapSys",
			value: MetricsValueToString(m.memStats.m.HeapSys),
		},
		{
			mType: shared.GaugeType,
			name:  "HeapReleased",
			value: MetricsValueToString(m.memStats.m.HeapReleased),
		},
		{
			mType: shared.GaugeType,
			name:  "HeapObjects",
			value: MetricsValueToString(m.memStats.m.HeapObjects),
		},
		{
			mType: shared.GaugeType,
			name:  "HeapInuse",
			value: MetricsValueToString(m.memStats.m.HeapInuse),
		},
		{
			mType: shared.GaugeType,
			name:  "HeapIdle",
			value: MetricsValueToString(m.memStats.m.HeapIdle),
		},
		{
			mType: shared.GaugeType,
			name:  "HeapAlloc",
			value: MetricsValueToString(m.memStats.m.HeapAlloc),
		},
		{
			mType: shared.GaugeType,
			name:  "GCSys",
			value: MetricsValueToString(m.memStats.m.GCSys),
		},
		{
			mType: shared.GaugeType,
			name:  "NumForcedGC",
			value: MetricsValueToString(m.memStats.m.NumForcedGC),
		},
		{
			mType: shared.GaugeType,
			name:  "NumGC",
			value: MetricsValueToString(m.memStats.m.NumGC),
		},
		{
			mType: shared.GaugeType,
			name:  "OtherSys",
			value: MetricsValueToString(m.memStats.m.OtherSys),
		},
		{
			mType: shared.GaugeType,
			name:  "PauseTotalNs",
			value: MetricsValueToString(m.memStats.m.PauseTotalNs),
		},
		{
			mType: shared.GaugeType,
			name:  "StackInuse",
			value: MetricsValueToString(m.memStats.m.StackInuse),
		},
		{
			mType: shared.GaugeType,
			name:  "StackSys",
			value: MetricsValueToString(m.memStats.m.StackSys),
		},
		{
			mType: shared.GaugeType,
			name:  "Sys",
			value: MetricsValueToString(m.memStats.m.Sys),
		},
		{
			mType: shared.GaugeType,
			name:  "TotalAlloc",
			value: MetricsValueToString(m.memStats.m.TotalAlloc),
		},
		{
			mType: shared.GaugeType,
			name:  "RandomValue",
			value: MetricsValueToString(m.memStats.RandomValue),
		},
		{
			mType: shared.CounterType,
			name:  "PollCount",
			value: MetricsValueToString(m.memStats.PollCount),
		},
	}
}

func (m RuntimeMetrics) Sync() {

	m.memStats.mutex.Lock()
	defer m.memStats.mutex.Unlock()

	runtime.ReadMemStats(m.memStats.m)
	m.memStats.PollCount += 1
	m.memStats.RandomValue = m.generateRandomFloat()
}

func (m RuntimeMetrics) generateRandomFloat() float64 {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return random.Float64()
}
