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
			APIToken:     "1",
			Monitor:      []string{"light", "temperatureMeasurement", "illuminanceMeasurement", "relativeHumidityMeasurement", "ultravioletIndex"},
			Period:       120,
			InfluxURL:    "http://localhost:8086",
			InfluxToken:  "token",
			InfluxOrg:    "org",
			InfluxBucket: "bucket",
			ValueMap:     map[string]map[string]float64{"switch": map[string]float64{"on": 1, "off": 0}},
		}, wantErr: false},
		{name: "second", file: "testdata/second.yaml", want: &Config{
			Monitor: []string{"light", "temperatureMeasurement", "illuminanceMeasurement", "relativeHumidityMeasurement", "ultravioletIndex"},
			SmartThings: SmartThingsConfig{
				Capabilities: monitor.MonitorCapabilities{
					monitor.MonitorCapability{Name: "temperatureMeasurement", Time: monitor.WallTime},
				},
			},
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
	influx, err := database.NewInfluxDBClient("http://url", "token", "org", "database")
	if err != nil {
		t.Errorf("Could not initialize influx %v", err)
	}
	type fields struct {
		APIToken     string
		Monitor      []string
		Period       int
		InfluxURL    string
		InfluxToken  string
		InfluxOrg    string
		InfluxBucket string
		ValueMap     monitor.ConversionMap
		SmartThings  SmartThingsConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *monitor.Monitor
	}{
		{
			name:   "defaults",
			fields: fields{},
			want:   monitor.New(),
		},
		{
			name:   "client",
			fields: fields{APIToken: "token"},
			want:   monitor.New(monitor.SetClient(smartthings.New("token"))),
		},
		{
			name: "multiple monitors",
			fields: fields{APIToken: "token", Monitor: []string{"a", "b", "c"},
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
			fields: fields{APIToken: "token", Monitor: []string{"a", "b", "c"},
				InfluxURL:    "http://url",
				InfluxToken:  "token",
				InfluxOrg:    "org",
				InfluxBucket: "bucket",
				Period:       360,
				ValueMap:     map[string]map[string]float64{"switch": {"on": 1, "off": 0}},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIToken:     tt.fields.APIToken,
				Monitor:      tt.fields.Monitor,
				Period:       tt.fields.Period,
				InfluxURL:    tt.fields.InfluxURL,
				InfluxToken:  tt.fields.InfluxToken,
				InfluxOrg:    tt.fields.InfluxOrg,
				InfluxBucket: tt.fields.InfluxBucket,
				ValueMap:     tt.fields.ValueMap,
				SmartThings:  tt.fields.SmartThings,
			}
			if got := c.InstantiateMonitor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.InstantiateMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}
