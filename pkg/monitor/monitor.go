package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/eargollo/smartthings-influx/pkg/smartthings"
)

type Monitor struct {
	period       time.Duration
	client       smartthings.Client
	recorder     Recorder
	lastUpdate   map[uuid.UUID]time.Time
	clock        Clock
	capabilities map[string]*MonitorCapability
	converter    ConversionMap
}

// New creates a new monitor that will add read data from the client
// and add to the recorder. If no recorder is passed as options
// stdout is used.
func New(opts ...MonitorOption) *Monitor {
	// Monitor with defaults
	mon := &Monitor{
		recorder: &StdOutRecorder{},
		clock:    &realClock{},
		period:   10 * time.Second,
	}

	mon.lastUpdate = make(map[uuid.UUID]time.Time)
	mon.capabilities = make(map[string]*MonitorCapability)

	for _, opt := range opts {
		opt(mon)
	}

	return mon
}

func (mon Monitor) Run() error {
	// Cheap trick not to sleep at the first round
	duration := time.Duration(0)

	for {
		// Cheap trick not to sleep at the first round
		time.Sleep(duration)
		duration = time.Duration(mon.period)
		// End of cheap trick

		dataPoints, err := mon.InspectDevices()
		if err != nil {
			log.Printf("ERROR: Could not gather devices data: %v", err)

			continue
		}

		// Using another map so we update the timestamp only when the record is serialized
		newLastUpdate := make(map[uuid.UUID]time.Time)

		if len(dataPoints) == 0 {
			log.Printf("ERROR: no devices with any of the capabilities: %s", strings.Join(mon.CapabilityNames(), ", "))

			continue
		}

		updateDataPoints := []DeviceDataPoint{}
		// Check which devices measurements were updated since last time
		for _, dp := range dataPoints {
			if mon.lastUpdate[dp.DeviceId] != dp.Timestamp {
				// Device updated, add to the update list
				updateDataPoints = append(updateDataPoints, dp)
			} else {
				log.Printf("No changes since last query for device %s[%s]. Skipping.", dp.Device, dp.DeviceId)
			}
			newLastUpdate[dp.DeviceId] = dp.Timestamp
		}

		if len(updateDataPoints) > 0 {
			err = mon.recorder.Add(updateDataPoints)
			if err != nil {
				log.Printf("Monitor got error writing point: %v", err)
			} else {
				log.Printf("Record saved %v", dataPoints)
				// Replace last update timestamps
				mon.lastUpdate = newLastUpdate
			}
		} else {
			log.Printf("No new data since last update")
		}
	}
}

func (mon Monitor) InspectDevices() ([]DeviceDataPoint, error) {
	currentTime := mon.clock.Now()
	dataPoints := []DeviceDataPoint{}

	if mon.client == nil {
		return dataPoints, fmt.Errorf("Can't connect to SmartThings, client not configured")
	}

	// List devices with metrics

	devices, err := mon.DevicesWithCapabilities()
	if err != nil {
		return dataPoints, fmt.Errorf("could not list devices %v", err)
	}

	if len(devices) == 0 {
		log.Printf("no devices with any of the metrics: %s", strings.Join(mon.CapabilityNames(), ", "))
		return dataPoints, nil
	}

	for i, dev := range devices {
		log.Printf("%d: Monitoring '%s' from device '%s' (%s)", i, dev.CapabilityId, dev.DeviceLabel, dev.DeviceId)

		// Get measurement
		status, err := mon.client.DeviceCapabilityStatus(dev.DeviceId, dev.ComponentId, dev.CapabilityId)
		if err != nil {
			log.Printf("ERROR: could not get metric status: %v", err)
			continue
		}

		for key, val := range status {
			if val.Value == nil {
				log.Printf("WARNING: Got nil metric value: %v", err)
				continue
			}

			// Get converted value
			convValue, err := mon.converter.Convert(key, val.Value)
			if err != nil {
				log.Printf("ERROR: could not convert to number %v", err)
				continue
			}

			readTime := val.Timestamp
			mc, ok := mon.capabilities[dev.CapabilityId]

			// log.Printf("Key %s Device ID %s Device %s Component %s Capability %s Value %v number value %f mc %s, ok %v config %v",
			// 	key,
			// 	dev.DeviceId,
			// 	dev.DeviceLabel,
			// 	dev.ComponentId,
			// 	dev.CapabilityId,
			// 	val,
			// 	convValue,
			// 	mc,
			// 	ok,
			// 	mon.capabilities,
			// )

			if ok {
				if mc.Time == WallTime {
					readTime = currentTime
					log.Printf("Device %s is set to wall time. Device read time %s replaced bt %s.",
						dev.DeviceLabel, val.Timestamp.String(), readTime.String())
				}
			}

			// Create point
			point := DeviceDataPoint{
				Key:        key,
				DeviceId:   dev.DeviceId,
				Device:     dev.DeviceLabel,
				Component:  dev.ComponentId,
				Capability: dev.CapabilityId,
				Unit:       val.Unit,
				Value:      convValue,
				Timestamp:  readTime,
			}

			dataPoints = append(dataPoints, point)
		}
	}

	return dataPoints, nil
}

type deviceWithCapability struct {
	DeviceId     uuid.UUID
	DeviceLabel  string
	ComponentId  string
	CapabilityId string
}

func (mon Monitor) DevicesWithCapabilities() ([]deviceWithCapability, error) {
	list := []deviceWithCapability{}

	devices, err := mon.client.Devices()
	if err != nil {
		return list, err
	}

	for _, d := range devices.Items {
		for _, comp := range d.Components {
			for _, cap := range comp.Capabilities {
				_, ok := mon.capabilities[cap.Id]
				if ok {
					// Capability is being monitored
					list = append(list, deviceWithCapability{DeviceId: d.DeviceId, DeviceLabel: d.Label, ComponentId: comp.Id, CapabilityId: cap.Id})
				}
			}
		}
	}

	return list, nil
}

func (mon Monitor) CapabilityNames() []string {
	names := []string{}

	for k := range mon.capabilities {
		names = append(names, k)
	}

	return names
}
