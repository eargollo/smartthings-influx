package config

import (
	"fmt"
	"os"

	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/spf13/viper"
)

type Config struct {
	APIToken       string                         `yaml:"apitoken"`
	Monitor        []string                       `yaml:"monitor"`
	Period         int                            `yaml:"period"`
	InfluxURL      string                         `yaml:"influxurl"`
	InfluxUser     string                         `yaml:"influxuser"`
	InfluxPassword string                         `yaml:"influxpasswword"`
	InfluxDatabase string                         `yaml:"influxdatabase"`
	ValueMap       smartthings.ConversionMap      `yaml:"valuemap,omitempty"`
	MonitorConfig  map[string]MonitorConfguration `yaml:"monitorconfig,omitempty"`
}

type TimeRead string

const (
	Sensor TimeRead = "sensor"
	Call   TimeRead = "call"
)

type MonitorConfguration struct {
	TimeSet TimeRead
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
