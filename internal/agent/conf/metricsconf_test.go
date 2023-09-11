package conf

import "testing"

func TestMetricsConfig_PushAddr(t *testing.T) {
	type fields struct {
		PushAddress string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "1",
			fields: fields{
				PushAddress: "addr",
			},
			want: "addr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MetricsConfig{
				PushAddress: tt.fields.PushAddress,
			}
			if got := c.PushAddr(); got != tt.want {
				t.Errorf("PushAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
