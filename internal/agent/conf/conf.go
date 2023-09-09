package conf

import "time"

type AgentConf struct {
	PollIntervalSec   int
	ReportIntervalSec int
}

func NewAgentConf() *AgentConf {
	return &AgentConf{
		PollIntervalSec:   2,
		ReportIntervalSec: 10,
	}
}

func (c AgentConf) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalSec) * time.Second
}

func (c AgentConf) ReportInterval() time.Duration {
	return time.Duration(c.ReportIntervalSec) * time.Second
}
