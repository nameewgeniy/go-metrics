package service

import (
	"bytes"
	"compress/gzip"
	"go-metrics/internal/models"
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

	m, err := models.NewMetricsFactory().
		MakeFromMapForUpdateMetrics(vars)

	if err != nil {
		return err
	}

	jsonPayload, err := m.MarshalJSON()
	if err != nil {
		return err
	}

	// Создание нового буфера для хранения сжатых данных
	var gzippedPayload bytes.Buffer
	gzw := gzip.NewWriter(&gzippedPayload)

	// Запись JSON-полезной нагрузки в сжатый буфер
	_, err = gzw.Write(jsonPayload)
	if err != nil {
		return err
	}

	err = gzw.Close()
	if err != nil {
		return err
	}

	// Создание нового Reader из сжатого буфера для отправки в запросе
	gzippedBody := bytes.NewReader(gzippedPayload.Bytes())

	u := &url.URL{
		Scheme: "http",
		Host:   s.cf.PushAddr(),
	}

	u.Path = path.Join("update") + "/"

	// Установка заголовка Content-Encoding для указания использования сжатия gzip
	req, err := http.NewRequest("POST", u.String(), gzippedBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса с сжатым телом
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
