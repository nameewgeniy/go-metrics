package service

import (
	"bytes"
	"go-metrics/internal/shared/metrics"
	"net/http"
	"net/url"
	"path"
)

type JSONSender struct {
	cf SenderConfig
}

func NewMetricJSONSender(cf SenderConfig) *JSONSender {
	return &JSONSender{
		cf: cf,
	}
}

func (s JSONSender) SendMemStatsMetric(metricType, metricName, metricValue string) error {

	vars := map[string]string{
		"type":  metricType,
		"name":  metricName,
		"value": metricValue,
	}

	m, err := metrics.NewMetricsFactory().
		MakeFromMapForUpdateMetrics(vars)

	if err != nil {
		return err
	}

	jsonPayload, err := m.MarshalJSON()

	if err != nil {
		return err
	}

	body := bytes.NewReader(jsonPayload)

	u := &url.URL{
		Scheme: "http",
		Host:   s.cf.PushAddr(),
	}

	u.Path = path.Join("update") + "/"

	resp, err := http.Post(u.String(), "application/json", body)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}