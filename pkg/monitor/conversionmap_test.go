package monitor

import (
	"testing"
)

func TestClient_Convert(t *testing.T) {
	type args struct {
		metric string
		value  any
	}
	tests := []struct {
		name          string
		conversionMap ConversionMap
		args          args
		want          float64
		wantErr       bool
	}{
		{
			name:          "float and no map",
			conversionMap: map[string]map[string]float64{},
			args: args{
				metric: "temperature",
				value:  8.25,
			},
			want:    8.25,
			wantErr: false,
		},
		{
			name:          "string and no map",
			conversionMap: map[string]map[string]float64{},
			args: args{
				metric: "light",
				value:  "off",
			},
			want:    0,
			wantErr: true,
		},
		{
			name:          "mapped string",
			conversionMap: map[string]map[string]float64{"light": {"on": 1.0}},
			args: args{
				metric: "light",
				value:  "on",
			},
			want:    1.0,
			wantErr: false,
		},
		{
			name:          "mapped map err",
			conversionMap: map[string]map[string]float64{"light": {"on": 1.0}},
			args: args{
				metric: "light",
				value:  map[string]string{"a": "b"},
			},
			want:    0,
			wantErr: true,
		},
		{
			name:          "case difference on metric bug #33",
			conversionMap: map[string]map[string]float64{"carbonmonoxide": {"clear": 0.0, "present": 1.0}},
			args: args{
				metric: "carbonMonoxide",
				value:  "present",
			},
			want:    1.0,
			wantErr: false,
		},
		{
			name:          "case difference on value bug #33",
			conversionMap: map[string]map[string]float64{"carbonmonoxide": {"clear": 0.0, "present": 1.0}},
			args: args{
				metric: "carbonMonoxide",
				value:  "Present",
			},
			want:    1.0,
			wantErr: false,
		},
		{
			name:          "case difference on value err if non existing",
			conversionMap: map[string]map[string]float64{"carbonmonoxide": {"clear": 0.0, "present": 1.0}},
			args: args{
				metric: "carbonMonoxide",
				value:  "notthere",
			},
			want:    0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.conversionMap.Convert(tt.args.metric, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ConvertValueToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.ConvertValueToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
