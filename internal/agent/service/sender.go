package service

import (
	"bytes"
	"fmt"
	"net/http"
)

type MSender struct {
}

func NewMetricSender() *MSender {
	return &MSender{}
}

func (s MSender) SendMemStatsMetric(addr, metricType, metricName string, metricValue any) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", addr, metricType, metricName, metricValue)

	resp, err := http.Post(url, "text/plain", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
