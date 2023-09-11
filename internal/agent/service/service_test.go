package service

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"github.com/nameewgeniy/go-metrics/internal/agent/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"runtime"
	"sync"
	"testing"
)

type MockSender struct {
	mock.Mock
}

func (s *MockSender) SendMemStatsMetric(addr, metricType, metricName string, metricValue any) error {
	args := s.Called(addr, metricType, metricName, metricValue)

	return args.Error(0)
}

func TestRuntimeMetrics_Sync(t *testing.T) {
	mockSender := new(MockSender)

	type fields struct {
		memStats *extMemStats
		mt       *metricsTracked
		cf       agent.MetricsConf
		s        MetricsSender
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "1",
			fields: fields{
				memStats: &extMemStats{
					m:           &runtime.MemStats{},
					PollCount:   0,
					RandomValue: 0,
					mutex:       sync.RWMutex{},
				},
				mt: &metricsTracked{
					gauges:   [26]string{"MSpanSys"},
					counters: nil,
				},
				cf: nil,
				s:  mockSender,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := RuntimeMetrics{
				memStats: tt.fields.memStats,
				mt:       tt.fields.mt,
				cf:       tt.fields.cf,
				s:        tt.fields.s,
			}
			m.Sync()
			assert.True(t, tt.fields.memStats.PollCount == 1)
			assert.True(t, tt.fields.memStats.RandomValue > 0)
			assert.True(t, tt.fields.memStats.m.MSpanSys > 0)
		})
	}
}

func TestRuntimeMetrics_Push(t *testing.T) {

	mockSender := new(MockSender)

	type fields struct {
		memStats *extMemStats
		mt       *metricsTracked
		cf       agent.MetricsConf
		s        MetricsSender
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "1",
			fields: fields{
				memStats: &extMemStats{
					m:           &runtime.MemStats{},
					PollCount:   0,
					RandomValue: 0,
					mutex:       sync.RWMutex{},
				},
				mt: &metricsTracked{
					gauges:   [26]string{"MSpanSys"},
					counters: nil,
				},
				cf: conf.NewMetricsConf(),
				s:  mockSender,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := RuntimeMetrics{
				memStats: tt.fields.memStats,
				mt:       tt.fields.mt,
				cf:       tt.fields.cf,
				s:        tt.fields.s,
			}
			m.Push()
		})
	}
}
