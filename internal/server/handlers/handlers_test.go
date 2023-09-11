package handlers

import (
	"net/http"
	"testing"
)

func TestHandlers_UpdateMetricsHande(t *testing.T) {
	type fields struct {
		m Metrics
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				m: tt.fields.m,
			}
			h.UpdateMetricsHande(tt.args.w, tt.args.r)
		})
	}
}
