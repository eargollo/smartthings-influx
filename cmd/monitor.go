/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/eargollo/smartthings-influx/internal/config"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Continuously monitor your SmartThings",
	Long: `Monitor your SmartThings non stop according to the configuration
	at .smartthings-influx.yaml. Runs at regular intervals and saves the measurements
	data in InfluxDB.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.Load(cfgFile)
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}

		// Monitor
		mon := config.InstantiateMonitor()

		err = mon.Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
