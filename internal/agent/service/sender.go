package service

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"go-metrics/internal/shared/metrics"
	"net/http"
	"net/url"
	"time"
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

	req, err := retryablehttp.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = time.Second * 5
	retryClient.CheckRetry = retryablehttp.DefaultRetryPolicy

	resp, err := retryClient.Do(req)
	if err != nil {
		return err
	}

	resp, err = retryClient.Do(req)

	return resp.Body.Close()
}
