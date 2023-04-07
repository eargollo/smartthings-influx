package monitor

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Recorder interface {
	Add([]DeviceDataPoint) error
}

type DeviceDataPoint struct {
	Key        string
	DeviceId   uuid.UUID
	Device     string
	Component  string
	Capability string
	Unit       string
	Value      float64
	Timestamp  time.Time
}

type StdOutRecorder struct{}

func (s *StdOutRecorder) Add(out []DeviceDataPoint) error {
	for i, dp := range out {
		fmt.Printf("%d, %s, %s, %s, %s, %s, %s, %f, %s\n",
			i,
			dp.Timestamp,
			dp.Key,
			dp.DeviceId,
			dp.Device,
			dp.Component,
			dp.Capability,
			dp.Value,
			dp.Unit,
		)
	}
	return nil
}
