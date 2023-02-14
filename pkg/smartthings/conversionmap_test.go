package smartthings

import (
	"testing"
)

func TestClient_ConvertValueToFloat(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.conversionMap.ConvertValueToFloat(tt.args.metric, tt.args.value)
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
