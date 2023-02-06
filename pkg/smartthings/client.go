package smartthings

type Client interface {
	Devices() (devices DevicesList, err error)
	DevicesWithCapabilities(capabilities []string) (list DevicesWithCapabilitiesResult, err error)
	// DeviceStatus(deviceID uuid.UUID) (status DeviceStatus, err error)
	// DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]CapabilityStatus, err error)
	CapabilityStatusToFloat(metric string, status CapabilityStatus) (float64, error)
}
