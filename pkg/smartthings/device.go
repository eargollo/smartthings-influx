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

type DevicesList struct {
	Items []Device `json:"items"`
}

type DevicesWithCapabilitiesResult struct {
	Items []DeviceWithCapability
}

type DeviceWithCapability struct {
	Device       Device
	ComponentId  string
	CapabilityId string
}

type DeviceCapabilitiesResult struct {
	ComponentId  string
	CapabilityId string
}
