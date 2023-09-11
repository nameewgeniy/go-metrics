package service

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"testing"
)

func TestRuntimeMetrics_Push(t *testing.T) {
	type fields struct {
		memStats *extMemStats
		mt       *metricsTracked
		cf       agent.MetricsConf
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := RuntimeMetrics{
				memStats: tt.fields.memStats,
				mt:       tt.fields.mt,
				cf:       tt.fields.cf,
			}
			m.Push()
		})
	}
}
