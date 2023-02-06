package database

import (
	"time"

	"github.com/google/uuid"
)

type Client interface {
	Save([]DeviceDataPoint) error
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
