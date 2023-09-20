package service

import (
	"net/http"
	"net/url"
	"path"
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

func (s MSender) SendMemStatsMetric(metricType, metricName, metricValue string) error {

	u := &url.URL{
		Scheme: "http",
		Host:   s.cf.PushAddr(),
	}

	u.Path = path.Join("update", metricType, metricName, metricValue)

	resp, err := http.Post(u.String(), "text/plain; charset=utf-8", nil)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
