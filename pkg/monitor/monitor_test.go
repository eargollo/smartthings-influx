package monitor

import (
	"reflect"
	"testing"
	"time"

	"github.com/eargollo/smartthings-influx/pkg/database"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockedTransport struct {
	mock.Mock
}

func (m *MockedTransport) Devices() (devices smartthings.DevicesList, err error) {
	args := m.Called()
	return args.Get(0).(smartthings.DevicesList), args.Error(1)
}

func (m *MockedTransport) DeviceStatus(deviceID uuid.UUID) (status smartthings.DeviceStatus, err error) {
	args := m.Called(deviceID)
	return args.Get(0).(smartthings.DeviceStatus), args.Error(1)
}

func (m *MockedTransport) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]smartthings.CapabilityStatus, err error) {
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

	cli := smartthings.Init(testObj, map[string]map[string]float64{})

	type fields struct {
		st       smartthings.Client
		metrics  []string
		interval int
	}

	tests := []struct {
		name    string
		fields  fields
		want    []database.DeviceDataPoint
		wantErr bool
	}{
		{
			name: "Keep original time",
			fields: fields{
				st:       cli,
				metrics:  []string{"temperatureMeasurement"},
				interval: 100,
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
			mon := New(tt.fields.st, tt.fields.metrics, tt.fields.interval)
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
