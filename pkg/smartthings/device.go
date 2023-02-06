package smartthings

import (
	"github.com/google/uuid"
)

type Device struct {
	DeviceId   uuid.UUID   `json:"deviceId"`
	Name       string      `json:"name"`
	Label      string      `json:"label"`
	Components []Component `json:"components"`
	Client     Transport
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

func (d *Device) Status() (DeviceStatus, error) {
	return d.Client.DeviceStatus(d.DeviceId)
}

func (d *Device) CapabilityStatus(componentId, capabilityID string) (map[string]CapabilityStatus, error) {
	return d.Client.DeviceCapabilityStatus(d.DeviceId, componentId, capabilityID)
}

type DeviceCapabilitiesResult struct {
	ComponentId  string
	CapabilityId string
}

func (d *Device) ListSelectedCapabilities(capabilities []string) []DeviceCapabilitiesResult {
	result := []DeviceCapabilitiesResult{}

	for _, comp := range d.Components {
		for _, cap := range comp.Capabilities {
			for _, m := range capabilities {
				if m == cap.Id {
					result = append(result, DeviceCapabilitiesResult{ComponentId: comp.Id, CapabilityId: cap.Id})
				}
			}
		}
	}

	return result
}
