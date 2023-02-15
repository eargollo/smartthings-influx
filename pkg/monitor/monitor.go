package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"

	"github.com/eargollo/smartthings-influx/internal/config"
	"github.com/eargollo/smartthings-influx/pkg/database"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
)

type Monitor struct {
	config     *config.Config
	stClient   smartthings.Client
	dbClient   database.Client
	lastUpdate map[uuid.UUID]time.Time
	clock      Clock
}

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func New(config *config.Config) (*Monitor, error) {
	mon := Monitor{config: config, clock: realClock{}}

	// SmartThings client if token is set
	if config.APIToken != "" {
		mon.stClient = smartthings.Init(smartthings.NewTransport(config.APIToken), config.ValueMap)
	}

	// dbClient
	if config.InfluxURL != "" {
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     viper.GetString("influxurl"),
			Username: viper.GetString("influxuser"),
			Password: viper.GetString("influxpassword"),
		})

		mon.dbClient = database.NewInfluxDBClient(c, config.InfluxDatabase)
		if err != nil {
			return &mon, err
		}
	}

	mon.lastUpdate = make(map[uuid.UUID]time.Time)

	return &mon, nil
}

func (mon *Monitor) SetClock(clock Clock) {
	mon.clock = clock
}

func (mon Monitor) Run() error {
	if mon.dbClient == nil {
		return fmt.Errorf("Can't monitor cause database is not set")
	}

	duration := time.Duration(0) // Cheap trick not to sleep at the first round

	for {
		// Cheap trick not to sleep at the first round
		time.Sleep(duration)
		duration = time.Duration(mon.config.Period) * time.Second
		// End of cheap trick

		dataPoints, err := mon.InspectDevices()
		if err != nil {
			log.Printf("ERROR: Could not gather devices data: %v", err)

			continue
		}

		// Using another map so we update the timestamp only when the record is serialized
		newLastUpdate := make(map[uuid.UUID]time.Time)

		if len(dataPoints) == 0 {
			log.Printf("ERROR: no devices with any of the metrics: %s", strings.Join(mon.config.Monitor, ", "))

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

	if mon.stClient == nil {
		return dataPoints, fmt.Errorf("Can't connect to SmartThings, client not configured")
	}

	// List devices with metrics
	devices, err := mon.stClient.DevicesWithCapabilities(mon.config.Monitor)
	if err != nil {
		return dataPoints, fmt.Errorf("could not list devices %v", err)
	}

	if len(devices.Items) == 0 {
		log.Printf("no devices with any of the metrics: %s", strings.Join(mon.config.Monitor, ", "))
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
			convValue, err := mon.stClient.CapabilityStatusToFloat(key, val)
			if err != nil {
				log.Printf("ERROR: could not convert to number %v", err)
				continue
			}

			log.Printf("Key is %s value %v number value %f", key, val, convValue)

			readTime := val.Timestamp
			mc, ok := mon.config.MonitorConfig[dev.CapabilityId]
			if ok {
				if mc.TimeSet == config.Call {
					readTime = mon.clock.Now()
				}
			}

			// Create point
			point := database.DeviceDataPoint{
				Key:        key,
				DeviceId:   dev.Device.DeviceId,
				Device:     dev.Device.Label,
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

func (mon *Monitor) SetTransport(transport smartthings.Transport) {
	mon.stClient = smartthings.Init(transport, mon.config.ValueMap)
}
