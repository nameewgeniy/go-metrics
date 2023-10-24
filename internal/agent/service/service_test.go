package service

import (
	"github.com/stretchr/testify/assert"
	"go-metrics/internal/shared/metrics"
	"testing"
	"time"
)

type MockMetricsSender struct {
	sentMetrics []metrics.Metrics
}

func (m *MockMetricsSender) SendMemStatsMetric(mt []metrics.Metrics) error {
	m.sentMetrics = mt
	return nil
}

func TestPush(t *testing.T) {
	mockSender := &MockMetricsSender{}
	runtimeMetrics := NewRuntimeMetrics(mockSender)

	runtimeMetrics.Push()

	// Allow some time for the goroutines to execute and send metrics
	time.Sleep(time.Second)

	expectedMetrics := runtimeMetrics.MetricsTracked()

	for _, v := range expectedMetrics {
		for _, i := range mockSender.sentMetrics {
			if v.ID == i.ID {
				assert.Equal(t, v.Value, i.Value, "Expected sent metrics value to match")
				assert.Equal(t, v.MType, i.MType, "Expected sent metrics type to match")
			}
		}
	}

}
