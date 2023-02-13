package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("APITOKEN", "1")

	tests := []struct {
		name    string
		file    string
		want    *Config
		wantErr bool
	}{
		{name: "empty", file: "", want: &Config{}, wantErr: false},
		{name: "first", file: "testdata/first.yaml", want: &Config{
			APIToken:       "1",
			Monitor:        []string{"light", "temperatureMeasurement", "illuminanceMeasurement", "relativeHumidityMeasurement", "ultravioletIndex"},
			Period:         120,
			InfluxURL:      "http://localhost:8086",
			InfluxUser:     "user",
			InfluxPassword: "password",
			InfluxDatabase: "database",
			ValueMap:       map[string]map[string]float64{"switch": map[string]float64{"on": 1, "off": 0}},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
