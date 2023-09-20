package conf

import "time"

type AgentConfig struct {
	PollIntervalSec   int
	ReportIntervalSec int
}

func NewAgentConf(pollIntervalSec, reportIntervalSec int) *AgentConfig {
	return &AgentConfig{
		PollIntervalSec:   pollIntervalSec,
		ReportIntervalSec: reportIntervalSec,
	}
}

func (c AgentConfig) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalSec) * time.Second
}

func (c AgentConfig) ReportInterval() time.Duration {
	return time.Duration(c.ReportIntervalSec) * time.Second
}
