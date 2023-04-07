package monitor

type ReadTime string

const (
	SensorTime ReadTime = "sensor"
	WallTime   ReadTime = "wall"
)

type MonitorCapability struct {
	Name string
	Time ReadTime
}

type MonitorCapabilities []MonitorCapability
