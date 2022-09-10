package smartthings

import (
	"fmt"
)

type DevicesList struct {
	Items  []Device `json:"items"`
	client *Client
}

type DevicesWithCapabilitiesResult struct {
	Items []DeviceWithCapability
}

type DeviceWithCapability struct {
	Device     Device
	Component  Component
	Capability Capability
	client     *Client
}

func (dev DeviceWithCapability) Status() (status map[string]CapabilityStatus, err error) {
	if dev.client == nil {
		return status, fmt.Errorf("invalid client: nil")
	}

	return cli.DeviceCapabilityStatus(dev.Device.DeviceId, dev.Component.Id, dev.Capability.Id)
}

func (list *DevicesList) DevicesWithCapabilities(capabilities []string) (devices DevicesWithCapabilitiesResult) {
	for _, d := range list.Items {
		for _, comp := range d.Components {
			for _, cap := range comp.Capabilities {
				for _, m := range capabilities {
					if m == cap.Id {
						devices.Items = append(devices.Items, DeviceWithCapability{Device: d, Component: comp, Capability: cap, client: list.client})
					}
				}
			}
		}
	}

	return
}
