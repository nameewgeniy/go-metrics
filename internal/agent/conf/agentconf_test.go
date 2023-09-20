package conf

import (
	"testing"
	"time"
)

func TestAgentConfig_PollInterval(t *testing.T) {
	type fields struct {
		PollIntervalSec   int
		ReportIntervalSec int
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "1",
			fields: fields{
				PollIntervalSec:   1,
				ReportIntervalSec: 0,
			},
			want: time.Duration(1) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := AgentConfig{
				PollIntervalSec:   tt.fields.PollIntervalSec,
				ReportIntervalSec: tt.fields.ReportIntervalSec,
			}
			if got := c.PollInterval(); got != tt.want {
				t.Errorf("PollInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentConfig_ReportInterval(t *testing.T) {
	type fields struct {
		PollIntervalSec   int
		ReportIntervalSec int
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "1",
			fields: fields{
				PollIntervalSec:   0,
				ReportIntervalSec: 1,
			},
			want: time.Duration(1) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := AgentConfig{
				PollIntervalSec:   tt.fields.PollIntervalSec,
				ReportIntervalSec: tt.fields.ReportIntervalSec,
			}
			if got := c.ReportInterval(); got != tt.want {
				t.Errorf("ReportInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
