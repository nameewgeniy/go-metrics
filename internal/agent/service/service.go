package service

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	counterType = "counter"
	gaugeType   = "gauge"
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
			mType: gaugeType,
			name:  "Alloc",
			value: MetricsValueToString(m.memStats.m.Alloc),
		},
		{
			mType: gaugeType,
			name:  "BuckHashSys",
			value: MetricsValueToString(m.memStats.m.BuckHashSys),
		},
		{
			mType: gaugeType,
			name:  "Frees",
			value: MetricsValueToString(m.memStats.m.Frees),
		},
		{
			mType: gaugeType,
			name:  "GCCPUFraction",
			value: MetricsValueToString(m.memStats.m.GCCPUFraction),
		},
		{
			mType: gaugeType,
			name:  "Mallocs",
			value: MetricsValueToString(m.memStats.m.Mallocs),
		},
		{
			mType: gaugeType,
			name:  "MSpanSys",
			value: MetricsValueToString(m.memStats.m.MSpanSys),
		},
		{
			mType: gaugeType,
			name:  "MSpanInuse",
			value: MetricsValueToString(m.memStats.m.MSpanInuse),
		},
		{
			mType: gaugeType,
			name:  "MCacheSys",
			value: MetricsValueToString(m.memStats.m.MCacheSys),
		},
		{
			mType: gaugeType,
			name:  "MCacheInuse",
			value: MetricsValueToString(m.memStats.m.MCacheInuse),
		},
		{
			mType: gaugeType,
			name:  "Lookups",
			value: MetricsValueToString(m.memStats.m.Lookups),
		},
		{
			mType: gaugeType,
			name:  "LastGC",
			value: MetricsValueToString(m.memStats.m.LastGC),
		},
		{
			mType: gaugeType,
			name:  "HeapSys",
			value: MetricsValueToString(m.memStats.m.HeapSys),
		},
		{
			mType: gaugeType,
			name:  "HeapReleased",
			value: MetricsValueToString(m.memStats.m.HeapReleased),
		},
		{
			mType: gaugeType,
			name:  "HeapObjects",
			value: MetricsValueToString(m.memStats.m.HeapObjects),
		},
		{
			mType: gaugeType,
			name:  "HeapInuse",
			value: MetricsValueToString(m.memStats.m.HeapInuse),
		},
		{
			mType: gaugeType,
			name:  "HeapIdle",
			value: MetricsValueToString(m.memStats.m.HeapIdle),
		},
		{
			mType: gaugeType,
			name:  "HeapAlloc",
			value: MetricsValueToString(m.memStats.m.HeapAlloc),
		},
		{
			mType: gaugeType,
			name:  "GCSys",
			value: MetricsValueToString(m.memStats.m.GCSys),
		},
		{
			mType: gaugeType,
			name:  "NumForcedGC",
			value: MetricsValueToString(m.memStats.m.NumForcedGC),
		},
		{
			mType: gaugeType,
			name:  "NumGC",
			value: MetricsValueToString(m.memStats.m.NumGC),
		},
		{
			mType: gaugeType,
			name:  "OtherSys",
			value: MetricsValueToString(m.memStats.m.OtherSys),
		},
		{
			mType: gaugeType,
			name:  "PauseTotalNs",
			value: MetricsValueToString(m.memStats.m.PauseTotalNs),
		},
		{
			mType: gaugeType,
			name:  "StackInuse",
			value: MetricsValueToString(m.memStats.m.StackInuse),
		},
		{
			mType: gaugeType,
			name:  "StackSys",
			value: MetricsValueToString(m.memStats.m.StackSys),
		},
		{
			mType: gaugeType,
			name:  "Sys",
			value: MetricsValueToString(m.memStats.m.Sys),
		},
		{
			mType: gaugeType,
			name:  "TotalAlloc",
			value: MetricsValueToString(m.memStats.m.TotalAlloc),
		},
		{
			mType: gaugeType,
			name:  "RandomValue",
			value: MetricsValueToString(m.memStats.RandomValue),
		},
		{
			mType: counterType,
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
