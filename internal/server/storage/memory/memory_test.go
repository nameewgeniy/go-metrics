package memory

import (
	"github.com/nameewgeniy/go-metrics/internal/server/storage"
	"testing"
)

func TestMemory_Save(t *testing.T) {
	type fields struct {
		s *MemStorage
	}
	type args struct {
		i storage.MetricsItem
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Memory{
				s: tt.fields.s,
			}
			if err := m.Save(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
