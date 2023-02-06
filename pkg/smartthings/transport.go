package smartthings

import "github.com/google/uuid"

type Transport interface {
	Devices() (devices DevicesList, err error)
	DeviceStatus(deviceID uuid.UUID) (status DeviceStatus, err error)
	DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]CapabilityStatus, err error)
}
