package smartthings

import (
	"github.com/google/uuid"
)

type Device struct {
	DeviceId   uuid.UUID   `json:"deviceId"`
	Name       string      `json:"name"`
	Label      string      `json:"label"`
	Components []Component `json:"components"`
}

type Component struct {
	Id           string       `json:"id"`
	Label        string       `json:"label"`
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

type DeviceStatus map[string]interface{}

func (d *Device) Status() (DeviceStatus, error) {
	return cli.DeviceStatus(d.DeviceId)
}
