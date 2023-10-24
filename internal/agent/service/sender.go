package service

import (
	"bytes"
	"encoding/json"
	"go-metrics/internal/shared/metrics"
	"net/http"
	"net/url"
)

type SenderConfig interface {
	PushAddr() string
}

type MSender struct {
	cf SenderConfig
}

func NewMetricSender(cf SenderConfig) *MSender {
	return &MSender{
		cf: cf,
	}
}

func (s MSender) SendMemStatsMetric(metrics []metrics.Metrics) error {

	u := &url.URL{
		Scheme: "http",
		Host:   s.cf.PushAddr(),
		Path:   "updates/",
	}

	jsonPayload, err := json.Marshal(metrics)

	if err != nil {
		return err
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
