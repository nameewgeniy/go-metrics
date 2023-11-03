package service

import (
	"go-metrics/internal/shared"
	"go-metrics/internal/shared/logger"
	"go-metrics/internal/shared/metrics"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type MetricsSender interface {
	SendMemStatsMetric([]metrics.Metrics) error
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

	if err := m.s.SendMemStatsMetric(m.MetricsTracked()); err != nil {
		logger.Log.Error(err.Error())
	}
}

func (m RuntimeMetrics) MetricsTracked() []metrics.Metrics {
	return []metrics.Metrics{
		{
			MType: shared.GaugeType,
			ID:    "Alloc",
			Value: func() *float64 {
				v := float64(m.memStats.m.Alloc)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "NextGC",
			Value: func() *float64 {
				v := float64(m.memStats.m.NextGC)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "BuckHashSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.BuckHashSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "Frees",
			Value: func() *float64 {
				v := float64(m.memStats.m.Frees)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "GCCPUFraction",
			Value: &m.memStats.m.GCCPUFraction,
		},
		{
			MType: shared.GaugeType,
			ID:    "Mallocs",
			Value: func() *float64 {
				v := float64(m.memStats.m.Mallocs)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "MSpanSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.MSpanSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "MSpanInuse",
			Value: func() *float64 {
				v := float64(m.memStats.m.MSpanInuse)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "MCacheSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.MCacheSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "MCacheInuse",
			Value: func() *float64 {
				v := float64(m.memStats.m.MCacheInuse)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "Lookups",
			Value: func() *float64 {
				v := float64(m.memStats.m.Lookups)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "LastGC",
			Value: func() *float64 {
				v := float64(m.memStats.m.LastGC)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "HeapSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.HeapSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "HeapReleased",
			Value: func() *float64 {
				v := float64(m.memStats.m.HeapReleased)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "HeapObjects",
			Value: func() *float64 {
				v := float64(m.memStats.m.HeapObjects)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "HeapInuse",
			Value: func() *float64 {
				v := float64(m.memStats.m.HeapInuse)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "HeapIdle",
			Value: func() *float64 {
				v := float64(m.memStats.m.HeapIdle)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "HeapAlloc",
			Value: func() *float64 {
				v := float64(m.memStats.m.HeapAlloc)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "GCSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.GCSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "NumForcedGC",
			Value: func() *float64 {
				v := float64(m.memStats.m.NumForcedGC)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "NumGC",
			Value: func() *float64 {
				v := float64(m.memStats.m.NumGC)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "OtherSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.OtherSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "PauseTotalNs",
			Value: func() *float64 {
				v := float64(m.memStats.m.PauseTotalNs)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "StackInuse",
			Value: func() *float64 {
				v := float64(m.memStats.m.StackInuse)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "StackSys",
			Value: func() *float64 {
				v := float64(m.memStats.m.StackSys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "Sys",
			Value: func() *float64 {
				v := float64(m.memStats.m.Sys)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "TotalAlloc",
			Value: func() *float64 {
				v := float64(m.memStats.m.TotalAlloc)
				return &v
			}(),
		},
		{
			MType: shared.GaugeType,
			ID:    "RandomValue",
			Value: &m.memStats.RandomValue,
		},
		{
			MType: shared.CounterType,
			ID:    "PollCount",
			Delta: func() *int64 {
				v := MetricsValueUintToInt(m.memStats.PollCount)
				return &v
			}(),
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
