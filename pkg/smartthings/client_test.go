package smartthings

import (
	"reflect"
	"testing"
)

func Test_initConverstionMap(t *testing.T) {
	type args struct {
		in0 map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string]float64
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				in0: map[string]interface{}{},
			},
			want:    map[string]map[string]float64{},
			wantErr: false,
		},
		{
			name: "on off string",
			args: args{
				in0: map[string]interface{}{"light": map[string]string{"on": "1", "off": "0"}},
			},
			want:    map[string]map[string]float64{},
			wantErr: true,
		},
		{
			name: "on off number",
			args: args{
				in0: map[string]interface{}{"light": map[string]interface{}{"on": 1, "off": 0}},
			},
			want:    map[string]map[string]float64{"light": {"on": 1.0, "off": 0}},
			wantErr: false,
		},
		{
			name: "not a map",
			args: args{
				in0: map[string]interface{}{"light": "abracadabra"},
			},
			want:    map[string]map[string]float64{},
			wantErr: true,
		},
		{
			name: "not a number",
			args: args{
				in0: map[string]interface{}{"light": map[string]string{"on": "1", "off": "off"}},
			},
			want:    map[string]map[string]float64{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConversionMap(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("initConverstionMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initConverstionMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ConvertValueToFloat(t *testing.T) {
	type fields struct {
		token         string
		conversionMap map[string]map[string]float64
	}
	type args struct {
		metric string
		value  any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:   "float and no map",
			fields: fields{"", map[string]map[string]float64{}},
			args: args{
				metric: "temperature",
				value:  8.25,
			},
			want:    8.25,
			wantErr: false,
		},
		{
			name:   "string and no map",
			fields: fields{"", map[string]map[string]float64{}},
			args: args{
				metric: "light",
				value:  "off",
			},
			want:    0,
			wantErr: true,
		},
		{
			name:   "mapped string",
			fields: fields{"", map[string]map[string]float64{"light": {"on": 1.0}}},
			args: args{
				metric: "light",
				value:  "on",
			},
			want:    1.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				token:         tt.fields.token,
				conversionMap: tt.fields.conversionMap,
			}
			got, err := c.ConvertValueToFloat(tt.args.metric, tt.args.value)
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
