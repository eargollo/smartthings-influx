package smartthings

var cli *STClient

type STClient struct {
	transport     Transport
	conversionMap ConversionMap
}

func Init(transport Transport, conversionMap map[string]map[string]float64) *STClient {
	cli = &STClient{transport: transport, conversionMap: conversionMap}

	return cli
}

func (c STClient) Devices() (DevicesList, error) {
	return c.transport.Devices()
}

func (c STClient) DevicesWithCapabilities(capabilities []string) (list DevicesWithCapabilitiesResult, err error) {
	devices, err := c.transport.Devices()
	if err != nil {
		return
	}

	for _, dev := range devices.Items {
		capList := dev.ListSelectedCapabilities(capabilities)
		for _, cap := range capList {
			list.Items = append(list.Items, DeviceWithCapability{Device: dev, ComponentId: cap.ComponentId, CapabilityId: cap.CapabilityId})
		}
	}

	return
}

func (c STClient) CapabilityStatusToFloat(metric string, status CapabilityStatus) (float64, error) {
	return c.conversionMap.ConvertValueToFloat(metric, status.Value)
}
