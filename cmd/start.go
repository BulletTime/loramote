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
	"fmt"

	"github.com/bullettime/rn2483"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rn2483.SetName(viper.GetString("device.serial"))
		rn2483.SetBaud(viper.GetInt("device.baud"))

		rn2483.Connect()
		defer rn2483.Disconnect()

		// Pause the lorawan layer
		// This returns a number of milliseconds that this layer stays paused
		// TODO start timer with the return value
		rn2483.MacPause()
		rn2483.RadioSetSyncWord(false)

		rn2483.RadioSetFrequency(uint32(viper.GetInt("gateway.frequency")))
		rn2483.RadioSetSpreadingFactor(uint8(viper.GetInt("gateway.sf")))

		rn2483.RadioSetBandWidth(uint16(viper.GetInt("gateway.bandwidth")))
		rn2483.RadioSetCodingRate(uint8(viper.GetInt("gateway.codingrate")))
		rn2483.RadioSetCrc(viper.GetBool("gateway.crc"))
		rn2483.RadioSetPower(int8(viper.GetInt("gateway.power")))

		// The watchdog timer could be disabled for continuous receiving
		// rn2483.RadioSetWatchDogTimer(0)

		// RadioRx test
		for {
			b := rn2483.RadioRxBlocking(0)

			if len(b) > 0 {
				fmt.Printf("Received packet: %q", b)
			}
		}

	},
}

func init() {
	gatewayCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
