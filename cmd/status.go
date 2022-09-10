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
	"encoding/json"
	"fmt"
	"log"

	"github.com/eargollo/smartthings-influx/pkg/smartthings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of a device or all devices",
	Long:  `Shows the status of a specific device or of all devices`,
	Run: func(cmd *cobra.Command, args []string) {
		convMap, err := smartthings.ParseConversionMap(viper.GetStringMap("valuemap"))
		if err != nil {
			log.Fatalf("Error initializing SmartThings client: %v", err)
		}

		smartthings.Init(viper.GetString("apitoken"), convMap)
		list, err := smartthings.Devices()
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			log.Printf("Listing status of all devices")
			for i, d := range list.Items {
				status, _ := d.Status()
				bs, _ := json.Marshal(status)

				fmt.Printf("%d: %s (%s): %v\n", i, d.Label, d.Name, string(bs))
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
