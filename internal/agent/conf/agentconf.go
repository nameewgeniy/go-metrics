package conf

import "time"

type AgentConfig struct {
	PollIntervalSec   int
	ReportIntervalSec int
	RateLimitNum      int
}

func NewAgentConf(pollIntervalSec, reportIntervalSec, rateLimit int) *AgentConfig {
	return &AgentConfig{
		PollIntervalSec:   pollIntervalSec,
		ReportIntervalSec: reportIntervalSec,
		RateLimitNum:      rateLimit,
	}
}

func (c AgentConfig) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalSec) * time.Second
}

func (c AgentConfig) ReportInterval() time.Duration {
	return time.Duration(c.ReportIntervalSec) * time.Second
}

func (c AgentConfig) RateLimit() int {
	return c.RateLimitNum
}
