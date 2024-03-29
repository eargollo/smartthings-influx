package database

import (
	"fmt"
	"log"

	"github.com/avast/retry-go"
	"github.com/eargollo/smartthings-influx/pkg/monitor"
	"github.com/influxdata/influxdb/client/v2"
	influxcli "github.com/influxdata/influxdb/client/v2"
)

type InfluxDB struct {
	client   influxcli.HTTPClient
	database string
}

func NewInfluxDBClient(url, user, password, database string) (*InfluxDB, error) {
	c, err := influxcli.NewHTTPClient(client.HTTPConfig{
		Addr:     url,
		Username: user,
		Password: password,
	})

	if err != nil {
		return nil, fmt.Errorf("could not instantiate http client for influx: %v", err)
	}

	return &InfluxDB{client: c, database: database}, nil
}

func (db InfluxDB) Add(datapoints []monitor.DeviceDataPoint) error {
	bp, err := influxcli.NewBatchPoints(influxcli.BatchPointsConfig{
		Database:  db.database,
		Precision: "s",
	})
	if err != nil {
		return fmt.Errorf("could not initialize points batch: %v", err)
	}

	for _, dp := range datapoints {
		// Create point
		point, err := influxcli.NewPoint(
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
		if err != nil {
			return fmt.Errorf("could not create influx point: %v", err)
		}

		bp.AddPoint(point)
	}

	if len(bp.Points()) > 0 {
		// Record points
		err := retry.Do(func() error {
			result := db.client.Write(bp)
			if result != nil {
				log.Printf("error writing point, will retry: %v", result)
			}
			return result
		})

		if err != nil {
			return fmt.Errorf("could not write set of points to InfluxDB: %v", err)
		}
	}

	return nil
}
