package conf

type SenderConfig struct {
	PushAddress string
}

func NewSenderConfig(pushAddress string) *SenderConfig {
	return &SenderConfig{
		PushAddress: pushAddress,
	}
}

func (c SenderConfig) PushAddr() string {
	return c.PushAddress
}
