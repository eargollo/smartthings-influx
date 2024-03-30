package database

import (
	"context"
	"fmt"

	"github.com/eargollo/smartthings-influx/pkg/monitor"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDBv2 struct {
	client    influxdb2.Client
	write_api api.WriteAPIBlocking
}

func NewInfluxDBv2Client(url string, token string, org string, bucket string) (*InfluxDBv2, error) {
	c := influxdb2.NewClient(url, token)
	if c == nil {
		return nil, fmt.Errorf("could not instantiate client for influx")
	}

	w := c.WriteAPIBlocking(org, bucket)
	if w == nil {
		return nil, fmt.Errorf("could not instantiate write api for influx")
	}

	return &InfluxDBv2{client: c, write_api: w}, nil
}

func (db InfluxDBv2) Add(datapoints []monitor.DeviceDataPoint) error {
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
		// write synchronously
		err := db.write_api.WritePoint(context.Background(), point)
		if err != nil {
			return err
		}
	}

	return nil
}
