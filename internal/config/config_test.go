package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/eargollo/smartthings-influx/pkg/database"
	"github.com/eargollo/smartthings-influx/pkg/monitor"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
)

func TestLoad(t *testing.T) {
	t.Setenv("APITOKEN", "1")

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
		{name: "second", file: "testdata/second.yaml", want: &Config{
			Monitor: []string{"light", "temperatureMeasurement", "illuminanceMeasurement", "relativeHumidityMeasurement", "ultravioletIndex"},
			SmartThings: SmartThingsConfig{
				Capabilities: monitor.MonitorCapabilities{
					monitor.MonitorCapability{Name: "temperatureMeasurement", Time: monitor.WallTime},
				},
			},
		}, wantErr: false},
		{name: "influx v2 base", file: "testdata/influxv2-base.yaml", want: &Config{
			APIToken: "1",
			Monitor:  []string{"light", "temperatureMeasurement", "illuminanceMeasurement", "relativeHumidityMeasurement", "ultravioletIndex"},
			Period:   120,
			Database: &DatabaseConfig{Type: "influxdbv2", URL: "http://localhost:8086", Token: "token", Org: "org", Bucket: "bucket"},
			ValueMap: map[string]map[string]float64{"switch": map[string]float64{"on": 1, "off": 0}},
		}, wantErr: false},
		{name: "influx v2 both", file: "testdata/influxv2-both.yaml", want: &Config{
			APIToken:       "1",
			Monitor:        []string{"light", "temperatureMeasurement", "illuminanceMeasurement", "relativeHumidityMeasurement", "ultravioletIndex"},
			Period:         120,
			InfluxURL:      "http://localhost:8086",
			InfluxUser:     "user",
			InfluxPassword: "password",
			InfluxDatabase: "database",
			Database:       &DatabaseConfig{Type: "influxdbv2", URL: "http://localhost:8086", Token: "token", Org: "org", Bucket: "bucket"},
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

func TestConfig_InstantiateMonitor(t *testing.T) {
	influx, err := database.NewInfluxDBClient("http://url", "user", "pass", "database")
	if err != nil {
		t.Errorf("Could not initialize influx %v", err)
	}
	tests := []struct {
		name   string
		config *Config
		want   *monitor.Monitor
	}{
		{
			name:   "defaults",
			config: &Config{},
			want:   monitor.New(),
		},
		{
			name:   "client",
			config: &Config{APIToken: "token"},
			want:   monitor.New(monitor.SetClient(smartthings.New("token"))),
		},
		{
			name: "multiple monitors",
			config: &Config{APIToken: "token", Monitor: []string{"a", "b", "c"},
				SmartThings: SmartThingsConfig{Capabilities: monitor.MonitorCapabilities{
					monitor.MonitorCapability{Name: "b", Time: monitor.WallTime}}},
			},
			want: monitor.New(
				monitor.SetClient(smartthings.New("token")),
				monitor.Capabilities(
					monitor.MonitorCapabilities{
						monitor.MonitorCapability{Name: "a", Time: monitor.SensorTime},
						monitor.MonitorCapability{Name: "b", Time: monitor.WallTime},
						monitor.MonitorCapability{Name: "c", Time: monitor.SensorTime},
					}),
			),
		},
		{
			name: "all in",
			config: &Config{APIToken: "token", Monitor: []string{"a", "b", "c"},
				InfluxURL:      "http://url",
				InfluxUser:     "user",
				InfluxPassword: "pass",
				InfluxDatabase: "database",
				Period:         360,
				ValueMap:       map[string]map[string]float64{"switch": {"on": 1, "off": 0}},
			},
			want: monitor.New(
				monitor.SetClient(smartthings.New("token")),
				monitor.Capabilities(
					monitor.MonitorCapabilities{
						monitor.MonitorCapability{Name: "a", Time: monitor.SensorTime},
						monitor.MonitorCapability{Name: "b", Time: monitor.SensorTime},
						monitor.MonitorCapability{Name: "c", Time: monitor.SensorTime},
					}),
				monitor.WithPeriod(6*time.Minute),
				monitor.SetRecorder(influx),
				monitor.WithConversion(map[string]map[string]float64{"switch": {"on": 1, "off": 0}}),
			),
		},
		{
			name: "all in",
			config: &Config{APIToken: "token", Monitor: []string{"a", "b", "c"},
				Database: &DatabaseConfig{Type: "influxdbv1", URL: "http://url", User: "user", Password: "pass", Database: "database"},
				Period:   360,
				ValueMap: map[string]map[string]float64{"switch": {"on": 1, "off": 0}},
			},
			want: monitor.New(
				monitor.SetClient(smartthings.New("token")),
				monitor.Capabilities(
					monitor.MonitorCapabilities{
						monitor.MonitorCapability{Name: "a", Time: monitor.SensorTime},
						monitor.MonitorCapability{Name: "b", Time: monitor.SensorTime},
						monitor.MonitorCapability{Name: "c", Time: monitor.SensorTime},
					}),
				monitor.WithPeriod(6*time.Minute),
				monitor.SetRecorder(influx),
				monitor.WithConversion(map[string]map[string]float64{"switch": {"on": 1, "off": 0}}),
			), // cant test influx 2 cause it is different at every instantiation
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.InstantiateMonitor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.InstantiateMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}
