package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/google/uuid"
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
	id         uuid.UUID
	device     string
	component  string
	capability string
	last       time.Time
}

func New(st *smartthings.Client, influx client.HTTPClient, database string, metrics []string, interval int) *Monitor {
	return &Monitor{st: st, influx: influx, database: database, metrics: metrics, interval: interval}
}

func (mon Monitor) Run() error {
	devices, err := mon.MonitoringDevices()
	if err != nil {
		return fmt.Errorf("could not list devices: %v", err)
	}
	if len(devices) == 0 {
		return fmt.Errorf("no devices with any of the metrics: %s", strings.Join(mon.metrics, ", "))
	}

	for _, dev := range devices {
		log.Printf("Monitoring %s from device %s (%s)", dev.capability, dev.device, dev.id.String())
	}

	for {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  mon.database,
			Precision: "s",
		})
		if err != nil {
			log.Printf("error creating batch points for influx: %v", err)
			time.Sleep(time.Duration(mon.interval) * time.Second)
			continue
		}

		for i, dev := range devices {
			// Get measurement
			status, err := mon.st.DeviceCapabilityStatus(dev.id, dev.component, dev.capability)
			if err != nil {
				log.Printf("could not get capability status: %v", err)
				continue
			}

			for key, val := range status {
				if val == nil {
					log.Printf("could not get value for capability status: %v", err)
					continue
				}

				log.Printf("Key is %s and value %v", key, val)
				inner, ok := val.(map[string]interface{})
				if !ok {
					log.Print("error, type was not interface")
					continue
				}
				// Get timestamp
				layout := "2006-01-02T15:04:05.000Z"
				str, ok := inner["timestamp"].(string)
				if !ok {
					log.Print("error, timestatmp was not a string")
					continue
				}
				t, err := time.Parse(layout, str)
				if err != nil {
					log.Printf("could not convert timestamp %s: %v", str, err)
					continue
				}

				if inner["value"] == nil {
					continue
				}

				if dev.last == t {
					continue
				}

				// Create point
				point, err := client.NewPoint(
					key,
					map[string]string{
						"device":     dev.device,
						"component":  dev.component,
						"capability": dev.capability,
					},
					map[string]interface{}{
						"value": inner["value"].(float64),
					},
					t,
				)
				if err != nil {
					log.Printf("could not create point: %v", err)
					time.Sleep(time.Duration(mon.interval) * time.Second)
					continue
				}

				bp.AddPoint(point)
				devices[i].last = t
			}
		}

		if len(bp.Points()) > 0 {
			// Record points
			err = mon.influx.Write(bp)
			if err != nil {
				log.Printf("Error writing point: %v", err)
			} else {
				log.Printf("Record saved %v", bp)
			}
		} else {
			log.Printf("No new read since last update")
		}

		time.Sleep(time.Duration(mon.interval) * time.Second)
	}

	// return nil
}

func (mon Monitor) MonitoringDevices() (devices []mondevice, err error) {
	list, err := mon.st.Devices()
	if err != nil {
		return
	}

	for _, d := range list.Items {
		for _, comp := range d.Components {
			for _, cap := range comp.Capabilities {
				for _, m := range mon.metrics {
					if m == cap.Id {
						devices = append(devices, mondevice{id: d.DeviceId, device: d.Label, component: comp.Id, capability: cap.Id})
					}
				}
			}
		}
	}

	return
}
