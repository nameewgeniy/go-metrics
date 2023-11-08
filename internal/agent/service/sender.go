package service

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"go-metrics/internal/shared"
	"go-metrics/internal/shared/metrics"
	"go-metrics/internal/shared/signature"
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

	body, err := s.makeBody(metrics)
	if err != nil {
		return err
	}

	hash := signature.Sign.Hash(body.Bytes())

	req, err := retryablehttp.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set(shared.HashHeaderName, hash)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = time.Second * 5
	retryClient.CheckRetry = retryablehttp.DefaultRetryPolicy

	resp, err := retryClient.Do(req)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func (s MSender) makeBody(metrics []metrics.Metrics) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	gz, _ := gzip.NewWriterLevel(buf, flate.BestCompression)

	if err := json.NewEncoder(gz).Encode(&metrics); err != nil {
		return nil, fmt.Errorf("encoding metrics: %w", err)
	}

	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("closing gzip: %w", err)
	}

	return buf, nil
}
