package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eargollo/smartthings-influx/pkg/database"
	"github.com/eargollo/smartthings-influx/pkg/monitor"
	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/spf13/viper"
)

type Config struct {
	APIToken       string                `yaml:"apitoken"`
	Monitor        []string              `yaml:"monitor"`
	Period         int                   `yaml:"period"`
	InfluxURL      string                `yaml:"influxurl"`
	InfluxUser     string                `yaml:"influxuser"`
	InfluxPassword string                `yaml:"influxpasswword"`
	InfluxDatabase string                `yaml:"influxdatabase"`
	ValueMap       monitor.ConversionMap `yaml:"valuemap,omitempty"`
	SmartThings    SmartThingsConfig     `yaml:"smartthings,omitempty"`
}

type SmartThingsConfig struct {
	Capabilities monitor.MonitorCapabilities `yaml:"capabilities,omitempty"`
}

func Load(cfgFile string) (*Config, error) {
	conf := &Config{}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return conf, fmt.Errorf("error getting home dir for default config file: %w", err)
		}

		// Search config in home directory with name ".smartthings-influx" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".smartthings-influx")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(conf)

	if err != nil {
		err = fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return conf, err
}

func (c *Config) InstantiateMonitor() *monitor.Monitor {
	parms := []monitor.MonitorOption{}

	if c.APIToken != "" {
		parms = append(parms, monitor.SetClient(smartthings.New(c.APIToken)))
	}

	if len(c.Monitor)+len(c.SmartThings.Capabilities) > 0 {
		caps := monitor.MonitorCapabilities{}
		for _, c := range c.Monitor {
			caps = append(caps, monitor.MonitorCapability{Name: c, Time: monitor.SensorTime})
		}

		caps = append(caps, c.SmartThings.Capabilities...)
		parms = append(parms, monitor.Capabilities(caps))
	}

	if c.InfluxDatabase != "" || c.InfluxPassword != "" || c.InfluxURL != "" || c.InfluxUser != "" {
		db, err := database.NewInfluxDBClient(c.InfluxURL, c.InfluxUser, c.InfluxPassword, c.InfluxDatabase)
		if err != nil {
			log.Fatalf("could not initialize influx: %v", err)
		}
		parms = append(parms, monitor.SetRecorder(db))
	}

	if c.Period != 0 {
		parms = append(parms, monitor.WithPeriod(time.Duration(c.Period)*time.Second))
	}

	if len(c.ValueMap) > 0 {
		parms = append(parms, monitor.WithConversion(c.ValueMap))
	}

	return monitor.New(parms...)
}
