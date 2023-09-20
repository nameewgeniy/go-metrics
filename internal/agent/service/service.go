package service

import (
	"fmt"
	"github.com/nameewgeniy/go-metrics/internal"
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

	fmt.Print("push")
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
			mType: internal.GaugeType,
			name:  "Alloc",
			value: MetricsValueToString(m.memStats.m.Alloc),
		},
		{
			mType: internal.GaugeType,
			name:  "BuckHashSys",
			value: MetricsValueToString(m.memStats.m.BuckHashSys),
		},
		{
			mType: internal.GaugeType,
			name:  "Frees",
			value: MetricsValueToString(m.memStats.m.Frees),
		},
		{
			mType: internal.GaugeType,
			name:  "GCCPUFraction",
			value: MetricsValueToString(m.memStats.m.GCCPUFraction),
		},
		{
			mType: internal.GaugeType,
			name:  "Mallocs",
			value: MetricsValueToString(m.memStats.m.Mallocs),
		},
		{
			mType: internal.GaugeType,
			name:  "MSpanSys",
			value: MetricsValueToString(m.memStats.m.MSpanSys),
		},
		{
			mType: internal.GaugeType,
			name:  "MSpanInuse",
			value: MetricsValueToString(m.memStats.m.MSpanInuse),
		},
		{
			mType: internal.GaugeType,
			name:  "MCacheSys",
			value: MetricsValueToString(m.memStats.m.MCacheSys),
		},
		{
			mType: internal.GaugeType,
			name:  "MCacheInuse",
			value: MetricsValueToString(m.memStats.m.MCacheInuse),
		},
		{
			mType: internal.GaugeType,
			name:  "Lookups",
			value: MetricsValueToString(m.memStats.m.Lookups),
		},
		{
			mType: internal.GaugeType,
			name:  "LastGC",
			value: MetricsValueToString(m.memStats.m.LastGC),
		},
		{
			mType: internal.GaugeType,
			name:  "HeapSys",
			value: MetricsValueToString(m.memStats.m.HeapSys),
		},
		{
			mType: internal.GaugeType,
			name:  "HeapReleased",
			value: MetricsValueToString(m.memStats.m.HeapReleased),
		},
		{
			mType: internal.GaugeType,
			name:  "HeapObjects",
			value: MetricsValueToString(m.memStats.m.HeapObjects),
		},
		{
			mType: internal.GaugeType,
			name:  "HeapInuse",
			value: MetricsValueToString(m.memStats.m.HeapInuse),
		},
		{
			mType: internal.GaugeType,
			name:  "HeapIdle",
			value: MetricsValueToString(m.memStats.m.HeapIdle),
		},
		{
			mType: internal.GaugeType,
			name:  "HeapAlloc",
			value: MetricsValueToString(m.memStats.m.HeapAlloc),
		},
		{
			mType: internal.GaugeType,
			name:  "GCSys",
			value: MetricsValueToString(m.memStats.m.GCSys),
		},
		{
			mType: internal.GaugeType,
			name:  "NumForcedGC",
			value: MetricsValueToString(m.memStats.m.NumForcedGC),
		},
		{
			mType: internal.GaugeType,
			name:  "NumGC",
			value: MetricsValueToString(m.memStats.m.NumGC),
		},
		{
			mType: internal.GaugeType,
			name:  "OtherSys",
			value: MetricsValueToString(m.memStats.m.OtherSys),
		},
		{
			mType: internal.GaugeType,
			name:  "PauseTotalNs",
			value: MetricsValueToString(m.memStats.m.PauseTotalNs),
		},
		{
			mType: internal.GaugeType,
			name:  "StackInuse",
			value: MetricsValueToString(m.memStats.m.StackInuse),
		},
		{
			mType: internal.GaugeType,
			name:  "StackSys",
			value: MetricsValueToString(m.memStats.m.StackSys),
		},
		{
			mType: internal.GaugeType,
			name:  "Sys",
			value: MetricsValueToString(m.memStats.m.Sys),
		},
		{
			mType: internal.GaugeType,
			name:  "TotalAlloc",
			value: MetricsValueToString(m.memStats.m.TotalAlloc),
		},
		{
			mType: internal.GaugeType,
			name:  "RandomValue",
			value: MetricsValueToString(m.memStats.RandomValue),
		},
		{
			mType: internal.CounterType,
			name:  "PollCount",
			value: MetricsValueToString(m.memStats.PollCount),
		},
	}
}

func (m RuntimeMetrics) Sync() {

	fmt.Print("sync")
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
