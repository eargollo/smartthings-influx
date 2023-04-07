package monitor

import (
	"time"

	"github.com/eargollo/smartthings-influx/pkg/smartthings"
)

type MonitorOption func(*Monitor)

func WithClock(clock Clock) MonitorOption {
	return func(m *Monitor) {
		m.clock = clock
	}
}

func SetRecorder(recorder Recorder) MonitorOption {
	return func(m *Monitor) {
		m.recorder = recorder
	}
}

func WithPeriod(period time.Duration) MonitorOption {
	return func(m *Monitor) {
		m.period = period
	}
}

func SetClient(client smartthings.Client) MonitorOption {
	return func(m *Monitor) {
		m.client = client
	}
}

// Capabilities set the capabilities being monitored
// if the array has the same capability more than once
// the last configuration present on the array is set
func Capabilities(capabilities MonitorCapabilities) MonitorOption {
	return func(m *Monitor) {
		for i := range capabilities {
			m.capabilities[capabilities[i].Name] = &capabilities[i]
		}
	}
}

func WithConversion(cmap ConversionMap) MonitorOption {
	return func(m *Monitor) {
		m.converter = cmap
	}
}

// func (mon *Monitor) SetTransport(transport smartthings.Transport) {
// 	mon.stClient = smartthings.Init(transport, mon.config.ValueMap)
// }
