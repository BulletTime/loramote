// Copyright © 2018 Sven Agneessens <sven.agneessens@gmail.com>
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/template"
	"github.com/apex/log"
)

const configTemplate = `# LoRaMote device configuration
mote:
  # Device Serial Address
  #
  # Example:
  # /dev/ttyUSB
  address: {{viper "mote.address"}}

  # Device Baud Rate
  #
  # Example:
  # 9600
  baud: {{viper "mote.baud"}}

  # Device Read Timeout
  #
  # The read timeout in milliseconds
  timeout: {{viper "mote.timeout"}}
`

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the loramote configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		funcMap := template.FuncMap{
			"viper": viper.GetString,
		}

		tmpl, err := template.New("config").Funcs(funcMap).Parse(configTemplate)
		if err != nil {
			log.Fatalf("[config] parsing: %s", err)
		}

		err = tmpl.Execute(os.Stdout, "config")
		if err != nil {
			log.Fatalf("[config] execution: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Mote defaults
	viper.SetDefault("mote.address", "/dev/cu.usbmodem14321")
	viper.SetDefault("mote.baud", 115200)
	viper.SetDefault("mote.timeout", 500)
}
