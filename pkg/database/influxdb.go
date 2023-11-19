package database

import (
	"fmt"

	"github.com/eargollo/smartthings-influx/pkg/monitor"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDB struct {
	client    influxdb2.Client
	write_api api.WriteAPI
}

func NewInfluxDBClient(url string, token string, org string, bucket string) (*InfluxDB, error) {
	c := influxdb2.NewClient(url, token)
	if c == nil {
		return nil, fmt.Errorf("could not instantiate client for influx")
	}

	w := c.WriteAPI(org, bucket)
	if w == nil {
		return nil, fmt.Errorf("could not instantiate write api for influx")
	}

	return &InfluxDB{client: c, write_api: w}, nil
}

func (db InfluxDB) Add(datapoints []monitor.DeviceDataPoint) error {
	for _, dp := range datapoints {
		// Create point
		point := influxdb2.NewPoint(
			dp.Key,
			map[string]string{
				"device":     dp.Device,
				"component":  dp.Component,
				"capability": dp.Capability,
				"unit":       dp.Unit,
			},
			map[string]interface{}{
				"value": dp.Value,
			},
			dp.Timestamp,
		)
		if point == nil {
			return fmt.Errorf("could not create influx point")
		}

		// Record point
		// write asynchronously
		db.write_api.WritePoint(point)
	}

	// Force all unwritten data to be sent
	db.write_api.Flush()

	return nil
}
