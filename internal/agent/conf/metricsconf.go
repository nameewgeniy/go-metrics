package conf

type MetricsConfig struct {
	PushAddress string
}

func NewMetricsConf() *MetricsConfig {
	return &MetricsConfig{
		PushAddress: "localhost:8080",
	}
}

func (c MetricsConfig) PushAddr() string {
	return c.PushAddress
}
