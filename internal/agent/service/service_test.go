package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MetricItem struct {
	mType string
	name  string
	value string
}

type MockMetricsSender struct {
	sentMetrics []MetricItem
}

func (m *MockMetricsSender) SendMemStatsMetric(metricType, name, value string) error {
	m.sentMetrics = append(m.sentMetrics, MetricItem{metricType, name, value})
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
			if v.Name == i.name {
				assert.Equal(t, v.Value, i.value, "Expected sent metrics value to match")
				assert.Equal(t, v.MType, i.mType, "Expected sent metrics type to match")
			}
		}
	}

}
