// Copyright © 2017 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"os"
	"strconv"

	"github.com/bullettime/logger"
	"github.com/bullettime/rn2483"
	"github.com/segmentio/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the gateway",
	Long:  `The gateway will be configured to listen on a specific frequency, spreading factor, ...`,
	Run: func(cmd *cobra.Command, args []string) {
		type yamlDevice struct {
			Serial string `yaml:"serial"`
			Baud   int    `yaml:"baud"`
		}

		type yamlGateway struct {
			SF         int  `yaml:"sf"`
			Frequency  int  `yaml:"frequency"`
			Bandwidth  int  `yaml:"bandwidth"`
			CodingRate int  `yaml:"codingrate"`
			CRC        bool `yaml:"crc"`
			Power      int  `yaml:"power"`
		}

		type yamlConfig struct {
			Device  yamlDevice  `yaml:"device"`
			Gateway yamlGateway `yaml:"gateway"`
		}

		var frequencies = []string{
			"868100000",
			"868300000",
			"868500000",
			"867100000",
			"867300000",
			"867500000",
			"867700000",
			"867900000",
		}

		var spreadingFactors = []string{
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
		}

		var bandwidths = []string{
			rn2483.BW1,
			rn2483.BW2,
			rn2483.BW3,
		}

		var codingRates = []string{
			rn2483.CR5,
			rn2483.CR6,
			rn2483.CR7,
			rn2483.CR8,
		}

		var powers = []string{
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
		}

		var (
			frequency  int
			sf         int
			bandwidth  int
			codingRate int
			crc        bool
			power      int
			err        error
		)

		sfs := spreadingFactors[prompt.Choose("Spreading factor", spreadingFactors)]
		sf, err = strconv.Atoi(sfs)
		if err != nil {
			logger.Warning.Println("Spreading factor can't be set")
		}

		frequencyS := frequencies[prompt.Choose("Frequency", frequencies)]
		frequency, err = strconv.Atoi(frequencyS)
		if err != nil {
			logger.Warning.Println("Frequency can't be set")
		}

		bandwidthS := bandwidths[prompt.Choose("Bandwidth", bandwidths)]
		bandwidth, err = strconv.Atoi(bandwidthS)
		if err != nil {
			logger.Warning.Println("Bandwidth can't be set")
		}

		codingRateS := codingRates[prompt.Choose("Coding Rate", codingRates)]
		codingRate, err = strconv.Atoi(codingRateS[2:])
		if err != nil {
			logger.Warning.Println("Coding rate can't be set")
		}

		crc = prompt.Confirm("Set CRC header?")

		powerS := powers[prompt.Choose("Power", powers)]
		power, err = strconv.Atoi(powerS)
		if err != nil {
			logger.Warning.Println("Power can't be set")
		}

		newConfig := &yamlConfig{
			Device: yamlDevice{
				Serial: viper.GetString("device.serial"),
				Baud:   viper.GetInt("device.baud"),
			},
			Gateway: yamlGateway{
				SF:         sf,
				Frequency:  frequency,
				Bandwidth:  bandwidth,
				CodingRate: codingRate,
				CRC:        crc,
				Power:      power,
			},
		}

		output, err := yaml.Marshal(newConfig)
		if err != nil {
			logger.Warning.Println("Failed to generate YAML:", err)
		}

		f, err := os.Create(viper.ConfigFileUsed())
		if err != nil {
			logger.Warning.Println("Failed to create file:", err)
		}
		defer f.Close()

		f.Write(output)
		logger.Debug.Println("New configuration file saved:", viper.ConfigFileUsed())

	},
}

func init() {
	gatewayCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
