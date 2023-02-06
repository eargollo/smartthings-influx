package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/eargollo/smartthings-influx/pkg/database"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
)

type Monitor struct {
	st         smartthings.Client
	dbClient   database.Client
	metrics    []string
	interval   int
	lastUpdate map[uuid.UUID]time.Time
}

func New(st smartthings.Client, dbClient database.Client, metrics []string, interval int) *Monitor {
	mon := Monitor{st: st, dbClient: dbClient, metrics: metrics, interval: interval}
	mon.lastUpdate = make(map[uuid.UUID]time.Time)
	return &mon
}

func (mon Monitor) Run() error {
	duration := time.Duration(0) // Cheap trick not to sleep at the first round

	for {
		// Cheap trick not to sleep at the first round
		time.Sleep(duration)
		duration = time.Duration(mon.interval) * time.Second
		// End of cheap trick

		dataPoints, err := mon.InspectDevices()
		if err != nil {
			log.Printf("ERROR: Could not gather devices data: %v", err)
			continue
		}

		// Using another map so we update the timestamp only when the record is serialized
		newLastUpdate := make(map[uuid.UUID]time.Time)

		if len(dataPoints) == 0 {
			log.Printf("ERROR: no devices with any of the metrics: %s", strings.Join(mon.metrics, ", "))
			continue
		}

		updateDataPoints := []database.DeviceDataPoint{}
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
			err = mon.dbClient.Save(updateDataPoints)
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

func (mon Monitor) InspectDevices() ([]database.DeviceDataPoint, error) {
	dataPoints := []database.DeviceDataPoint{}

	// List devices with metrics
	devices, err := mon.st.DevicesWithCapabilities(mon.metrics)
	if err != nil {
		return dataPoints, fmt.Errorf("could not list devices %v", err)
	}

	if len(devices.Items) == 0 {
		log.Printf("no devices with any of the metrics: %s", strings.Join(mon.metrics, ", "))
		return dataPoints, nil
	}

	for i, dev := range devices.Items {
		log.Printf("%d: Monitoring '%s' from device '%s' (%s)", i, dev.CapabilityId, dev.Device.Label, dev.Device.DeviceId)
		// Get measurement
		status, err := dev.Device.CapabilityStatus(dev.ComponentId, dev.CapabilityId)
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
			convValue, err := mon.st.CapabilityStatusToFloat(key, val)
			if err != nil {
				log.Printf("ERROR: could not convert to number %v", err)
				continue
			}

			log.Printf("Key is %s value %v number value %f", key, val, convValue)

			// Create point
			point := database.DeviceDataPoint{
				Key:        key,
				DeviceId:   dev.Device.DeviceId,
				Device:     dev.Device.Label,
				Component:  dev.ComponentId,
				Capability: dev.CapabilityId,
				Unit:       val.Unit,
				Value:      convValue,
				Timestamp:  val.Timestamp,
			}

			dataPoints = append(dataPoints, point)
		}

	}

	return dataPoints, nil
}
