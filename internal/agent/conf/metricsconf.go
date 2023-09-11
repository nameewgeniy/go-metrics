package conf

type MetricsConfig struct {
	PushAddress string
}

func NewMetricsConf(pushAddress string) *MetricsConfig {
	return &MetricsConfig{
		PushAddress: pushAddress,
	}
}

func (c MetricsConfig) PushAddr() string {
	return c.PushAddress
}
