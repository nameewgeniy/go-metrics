package service

import "testing"

func TestMetrics_Update(t *testing.T) {
	type fields struct {
		s Storage
	}
	type args struct {
		mType  string
		mName  string
		mValue string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metrics{
				s: tt.fields.s,
			}
			if err := m.Update(tt.args.mType, tt.args.mName, tt.args.mValue); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
