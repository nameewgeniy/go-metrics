package conf

import "time"

type AgentConfig struct {
	PollIntervalSec   int
	ReportIntervalSec int
}

func NewAgentConf() *AgentConfig {
	return &AgentConfig{
		PollIntervalSec:   2,
		ReportIntervalSec: 10,
	}
}

func (c AgentConfig) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalSec) * time.Second
}

func (c AgentConfig) ReportInterval() time.Duration {
	return time.Duration(c.ReportIntervalSec) * time.Second
}
