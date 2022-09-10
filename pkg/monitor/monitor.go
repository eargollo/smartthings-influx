package monitor

import (
	"log"
	"strings"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/google/uuid"

	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/influxdata/influxdb/client/v2"
)

type Monitor struct {
	st       *smartthings.Client
	influx   client.HTTPClient
	database string
	metrics  []string
	interval int
}

type mondevice struct {
	device     smartthings.Device
	component  smartthings.Component
	capability string
	last       time.Time
}

func New(st *smartthings.Client, influx client.HTTPClient, database string, metrics []string, interval int) *Monitor {
	return &Monitor{st: st, influx: influx, database: database, metrics: metrics, interval: interval}
}

func (mon Monitor) Run() error {
	duration := time.Duration(0) // Cheap trick not to sleep at the first round

	lastUpdate := make(map[uuid.UUID]time.Time)

	for {
		// Cheap trick not to sleep at the first round
		time.Sleep(duration)
		duration = time.Duration(mon.interval) * time.Second
		// End of cheap trick

		// Using another map so we update the timestamp only when the record is serialized
		newLastUpdate := make(map[uuid.UUID]time.Time)

		// List devices with metrics
		devices, err := mon.st.DevicesWithCapabilities(mon.metrics)
		if err != nil {
			log.Printf("ERROR: could not list devices: %v", err)
			continue
		}
		if len(devices.Items) == 0 {
			log.Printf("ERROR: no devices with any of the metrics: %s", strings.Join(mon.metrics, ", "))
			continue
		}

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  mon.database,
			Precision: "s",
		})
		if err != nil {
			log.Printf("error creating batch points for influx: %v", err)
			time.Sleep(time.Duration(mon.interval) * time.Second)
			continue
		}

		for i, dev := range devices.Items {
			log.Printf("%d: Monitoring '%s' from device '%s' (%s)", i, dev.Capability.Id, dev.Device.Label, dev.Device.DeviceId)
			// Get measurement
			status, err := dev.Status()
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
				convValue, err := val.FloatValue(key)
				if err != nil {
					log.Printf("ERROR: could not convert to number %v", err)
					continue
				}

				log.Printf("Key is %s value %v number value %f", key, val, convValue)

				if lastUpdate[dev.Device.DeviceId] == val.Timestamp {
					log.Printf("No changes since last query. Skipping.")
					continue
				}

				// Create point
				point, err := client.NewPoint(
					key,
					map[string]string{
						"device":     dev.Device.Label,
						"component":  dev.Component.Id,
						"capability": dev.Capability.Id,
						"unit":       val.Unit,
					},
					map[string]interface{}{
						"value": convValue,
					},
					val.Timestamp,
				)
				if err != nil {
					log.Printf("could not create point: %v", err)
					time.Sleep(time.Duration(mon.interval) * time.Second)
					continue
				}

				bp.AddPoint(point)
				newLastUpdate[dev.Device.DeviceId] = val.Timestamp
			}
		}

		if len(bp.Points()) > 0 {
			// Record points
			err := retry.Do(func() error {
				result := mon.influx.Write(bp)
				if result != nil {
					log.Printf("Error writing point: %v", result)
				}
				return result
			})
			if err != nil {
				log.Printf("Error writing point: %v", err)
			} else {
				log.Printf("Record saved %v", bp)
				lastUpdate = newLastUpdate
			}
		} else {
			log.Printf("No new read since last update")
		}

	}
}
