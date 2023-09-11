package service

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

const gaugeType = "gauge"
const counterType = "counter"

// Метрки, которые мы будем отправлять из пакета runtime
type metricsTracked struct {
	gauges   [26]string
	counters []string
}

type extMemStats struct {
	m           *runtime.MemStats
	PollCount   uint64
	RandomValue float64
	mutex       sync.RWMutex
}

type MetricsSender interface {
	SendMemStatsMetric(addr, metricType, metricName string, metricValue any) error
}

type RuntimeMetrics struct {
	memStats *extMemStats
	mt       *metricsTracked
	cf       agent.MetricsConf
	s        MetricsSender
}

func NewRuntimeMetrics(cf agent.MetricsConf, sender MetricsSender) *RuntimeMetrics {
	return &RuntimeMetrics{
		memStats: &extMemStats{
			m: &runtime.MemStats{},
		},
		s: sender,
		mt: &metricsTracked{
			gauges: [26]string{
				"Alloc",
				"BuckHashSys",
				"Frees",
				"GCCPUFraction",
				"Mallocs",
				"MSpanSys",
				"MSpanInuse",
				"MCacheSys",
				"MCacheInuse",
				"Lookups",
				"LastGC",
				"HeapSys",
				"HeapReleased",
				"HeapObjects",
				"HeapInuse",
				"HeapIdle",
				"HeapAlloc",
				"GCSys",
				"NumForcedGC",
				"NumGC",
				"OtherSys",
				"PauseTotalNs",
				"StackInuse",
				"StackSys",
				"Sys",
				"TotalAlloc",
			},
			counters: nil,
		},
		cf: cf,
	}
}

func (m RuntimeMetrics) Push() {

	m.memStats.mutex.RLock()
	defer m.memStats.mutex.RUnlock()

	for _, v := range m.mt.gauges {
		name := v

		// Получаем значение свойства по имени
		fieldValue := reflect.ValueOf(*m.memStats.m).FieldByName(name)

		// Проверяем, что свойство с таким именем существует
		if fieldValue.IsValid() {
			go func() {
				err := m.s.SendMemStatsMetric(m.cf.PushAddr(), gaugeType, name, fieldValue)
				if err != nil {
					log.Print(err)
				}
			}()
		}
	}

	// Отправялем кастомные метрики не из пакета runtime
	go func() {
		err := m.s.SendMemStatsMetric(m.cf.PushAddr(), gaugeType, "RandomValue", m.memStats.RandomValue)
		if err != nil {
			log.Print(err)
		}
	}()

	go func() {
		err := m.s.SendMemStatsMetric(m.cf.PushAddr(), counterType, "PollCount", m.memStats.PollCount)
		if err != nil {
			log.Print(err)
		}
	}()
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
