package conf

import "testing"

func TestServerConf_Addr(t *testing.T) {
	type fields struct {
		addr string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "1",
			fields: fields{
				addr: "8080",
			},
			want: "8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ServerConf{
				addr: tt.fields.addr,
			}
			if got := c.Addr(); got != tt.want {
				t.Errorf("Addr() = %v, want %v", got, tt.want)
			}
		})
	}
}
