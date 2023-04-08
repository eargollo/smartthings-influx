package monitor_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/eargollo/smartthings-influx/pkg/monitor"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockedSTClient struct {
	mock.Mock
}

func (m *MockedSTClient) Devices() (smartthings.DevicesList, error) {
	args := m.Called()
	return args.Get(0).(smartthings.DevicesList), args.Error(1)
}

func (m *MockedSTClient) DeviceStatus(deviceID uuid.UUID) (smartthings.DeviceStatus, error) {
	args := m.Called(deviceID)
	return args.Get(0).(smartthings.DeviceStatus), args.Error(1)
}

func (m *MockedSTClient) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (map[string]smartthings.CapabilityStatus, error) {
	args := m.Called(deviceID, componentId, capabilityId)
	return args.Get(0).(map[string]smartthings.CapabilityStatus), args.Error(1)
}

type MockedClock struct {
	mock.Mock
}

func (m *MockedClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func TestMonitor_InspectDevices(t *testing.T) {
	id1 := uuid.New()
	ts, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	tempValue := float64(21)

	testObj := new(MockedSTClient)

	testObj.On("Devices").Return(
		smartthings.DevicesList{
			Items: []smartthings.Device{
				smartthings.Device{
					DeviceId: id1,
					Name:     "Mocked Device",
					Label:    "Mocked Device",
					Components: []smartthings.Component{
						smartthings.Component{
							Id:    "main",
							Label: "main",
							Capabilities: []smartthings.Capability{
								smartthings.Capability{
									Id:      "temperatureMeasurement",
									Version: 1,
								},
							},
						},
					},
				},
			},
		},
		nil,
	)

	testObj.On("DeviceCapabilityStatus", id1, "main", "temperatureMeasurement").Return(
		map[string]smartthings.CapabilityStatus{
			"temperature": smartthings.CapabilityStatus{
				Timestamp: ts,
				Unit:      "C",
				Value:     tempValue,
			},
		},
		nil,
	)

	test2Obj := new(MockedSTClient)

	test2Obj.On("Devices").Return(
		smartthings.DevicesList{
			Items: []smartthings.Device{
				{
					DeviceId: id1,
					Name:     "CarbonMonoxideDetector",
					Label:    "Smoke/CO Garage",
					Components: []smartthings.Component{
						{
							Id:    "carbonMonoxideDetector",
							Label: "carbonMonoxideDetector",
							Capabilities: []smartthings.Capability{
								{
									Id:      "carbonMonoxideDetector",
									Version: 1,
								},
							},
						},
					},
				},
			},
		},
		nil,
	)

	test2Obj.On("DeviceCapabilityStatus", id1, "carbonMonoxideDetector", "carbonMonoxideDetector").Return(
		map[string]smartthings.CapabilityStatus{
			"carbonMonoxide": {
				Timestamp: ts,
				Value:     "clear",
			},
		},
		nil,
	)

	tests := []struct {
		name    string
		mon     *monitor.Monitor
		want    []monitor.DeviceDataPoint
		wantErr bool
	}{
		{
			name: "Keep original time",
			mon: monitor.New(
				monitor.SetClient(testObj),
				monitor.WithPeriod(100*time.Second),
				monitor.Capabilities(monitor.MonitorCapabilities{
					monitor.MonitorCapability{Name: "temperatureMeasurement", Time: monitor.SensorTime},
				}),
			),
			want: []monitor.DeviceDataPoint{
				monitor.DeviceDataPoint{
					Key:        "temperature",
					DeviceId:   id1,
					Device:     "Mocked Device",
					Component:  "main",
					Capability: "temperatureMeasurement",
					Value:      tempValue,
					Unit:       "C",
					Timestamp:  ts,
				},
			},
			wantErr: false,
		},
		{
			name: "Conversion as in bug #33",
			mon: monitor.New(
				monitor.SetClient(test2Obj),
				monitor.WithPeriod(100*time.Second),
				monitor.WithConversion(monitor.ConversionMap{
					"carbonmonoxide": map[string]float64{"clear": 0.0},
				}),
				monitor.Capabilities(monitor.MonitorCapabilities{
					monitor.MonitorCapability{Name: "carbonMonoxideDetector", Time: monitor.SensorTime},
				}),
			),
			want: []monitor.DeviceDataPoint{
				monitor.DeviceDataPoint{
					Key:        "carbonMonoxide",
					DeviceId:   id1,
					Device:     "Smoke/CO Garage",
					Component:  "carbonMonoxideDetector",
					Capability: "carbonMonoxideDetector",
					Value:      0.0,
					Timestamp:  ts,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mon.InspectDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Monitor.InspectDevices() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Monitor.InspectDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonitor_CallTime(t *testing.T) {
	id1 := uuid.New()
	sensorTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	readTime, _ := time.Parse(time.RFC3339, "2023-01-01T10:00:00Z")
	tempValue := float64(21)

	testObj := new(MockedSTClient)
	clockObj := new(MockedClock)

	clockObj.On("Now").Return(readTime)

	testObj.On("Devices").Return(
		smartthings.DevicesList{
			Items: []smartthings.Device{
				smartthings.Device{
					DeviceId: id1,
					Name:     "Mocked Device",
					Label:    "Mocked Device",
					Components: []smartthings.Component{
						smartthings.Component{
							Id:    "main",
							Label: "main",
							Capabilities: []smartthings.Capability{
								smartthings.Capability{
									Id:      "temperatureMeasurement",
									Version: 1,
								},
							},
						},
					},
				},
			},
		},
		nil,
	)

	testObj.On("DeviceCapabilityStatus", id1, "main", "temperatureMeasurement").Return(
		map[string]smartthings.CapabilityStatus{
			"temperature": smartthings.CapabilityStatus{
				Timestamp: sensorTime,
				Unit:      "C",
				Value:     tempValue,
			},
		},
		nil,
	)

	tests := []struct {
		name string
		mon  *monitor.Monitor
		// config  config.Config
		want    []monitor.DeviceDataPoint
		wantErr bool
	}{
		{
			name: "Keep original time",
			mon: monitor.New(
				monitor.SetClient(testObj),
				monitor.Capabilities(
					monitor.MonitorCapabilities{
						monitor.MonitorCapability{
							Name: "temperatureMeasurement",
							Time: monitor.WallTime,
						},
					},
				),
				monitor.WithClock(clockObj),
				monitor.WithPeriod(100*time.Second),
			),
			want: []monitor.DeviceDataPoint{
				{
					Key:        "temperature",
					DeviceId:   id1,
					Device:     "Mocked Device",
					Component:  "main",
					Capability: "temperatureMeasurement",
					Value:      tempValue,
					Unit:       "C",
					Timestamp:  readTime,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mon.InspectDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Monitor.InspectDevices() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Monitor.InspectDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}
