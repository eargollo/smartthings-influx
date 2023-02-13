package monitor_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/eargollo/smartthings-influx/internal/config"
	"github.com/eargollo/smartthings-influx/pkg/database"
	"github.com/eargollo/smartthings-influx/pkg/monitor"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedTransport struct {
	mock.Mock
}

func (m *MockedTransport) Devices() (smartthings.DevicesList, error) {
	args := m.Called()
	return args.Get(0).(smartthings.DevicesList), args.Error(1)
}

func (m *MockedTransport) DeviceStatus(deviceID uuid.UUID) (smartthings.DeviceStatus, error) {
	args := m.Called(deviceID)
	return args.Get(0).(smartthings.DeviceStatus), args.Error(1)
}

func (m *MockedTransport) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (map[string]smartthings.CapabilityStatus, error) {
	args := m.Called(deviceID, componentId, capabilityId)
	return args.Get(0).(map[string]smartthings.CapabilityStatus), args.Error(1)
}

func TestMonitor_InspectDevices(t *testing.T) {
	id1 := uuid.New()
	ts, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	tempValue := float64(21)

	testObj := new(MockedTransport)

	testObj.On("Devices").Return(
		smartthings.DevicesList{
			Items: []smartthings.Device{
				smartthings.Device{
					Client:   testObj,
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

	tests := []struct {
		name    string
		config  config.Config
		want    []database.DeviceDataPoint
		wantErr bool
	}{
		{
			name: "Keep original time",
			config: config.Config{
				Monitor: []string{"temperatureMeasurement"},
				Period:  100,
			},
			want: []database.DeviceDataPoint{
				database.DeviceDataPoint{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mon, err := monitor.New(&tt.config)
			assert.NoError(t, err)

			mon.SetTransport(testObj)

			got, err := mon.InspectDevices()
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
