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
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/bullettime/rn2483"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	latitude  float32
	longitude float32
)

// ddrCmd represents the ddr command
var ddrCmd = &cobra.Command{
	Use:   "ddr",
	Short: "Test DDR algorithm",
	Run: func(cmd *cobra.Command, args []string) {
		TestDDR()
	},
}

func init() {
	RootCmd.AddCommand(ddrCmd)

	ddrCmd.Flags().Uint8VarP(&dataRate, "datarate", "r", 0, "set data rate (default 0)")
	ddrCmd.Flags().Uint8VarP(&timeout, "timeout", "t", 10, "set timeout in minutes")

	ddrCmd.Flags().Float32Var(&latitude, "lat", 0, "set the latitude")
	ddrCmd.MarkFlagRequired("lat")
	ddrCmd.Flags().Float32Var(&longitude, "lon", 0, "set the longitude")
	ddrCmd.MarkFlagRequired("lon")
}

func TestDDR() {
	// setup loggers
	rn2483.ERROR = myLogger{level: "ERROR"}
	rn2483.WARN = myLogger{level: "WARN"}
	rn2483.DEBUG = myLogger{level: "DEBUG"}

	// initialize device
	rn2483.SetName(viper.GetString("mote.address"))
	rn2483.SetBaud(viper.GetInt("mote.baud"))
	rn2483.SetTimeout(time.Millisecond * time.Duration(viper.GetInt("mote.timeout")))

	// connect
	rn2483.Connect()
	defer rn2483.Disconnect()

	fmt.Println(rn2483.Version())

	// setup mac settings
	rn2483.MacReset(868)
	rn2483.MacSetDeviceEUI(viper.GetString("lora.deveui"))
	rn2483.MacSetApplicationEUI(viper.GetString("lora.appeui"))
	rn2483.MacSetApplicationKey(viper.GetString("lora.appkey"))
	rn2483.MacSetDataRate(dataRate)
	rn2483.MacSetPowerIndex(rn2483.DBm14)
	rn2483.MacSetADR(false)

	log.WithFields(log.Fields{
		"data rate": rn2483.MacGetDataRate(),
		"power":     rn2483.MacGetPowerIndex(),
		"adr":       rn2483.MacGetADR(),
	}).Info("mac settings configured")

	// join the network
	if !rn2483.MacJoin(rn2483.OTAA) {
		log.Fatal("could not join the network")
	}

	log.Info("connected to the network")

	// setup duty cycles - dingnet has 7 channels on
	nbChannels := 7
	dcycle := 100 / float32(nbChannels)
	for i := 0; i <= nbChannels; i++ {
		rn2483.MacSetChannelDutyCycle(uint8(i), dcycle)
	}

	log.WithField("duty cycle", dcycle).Info("new duty cycles configured")

	// setup DDR call
	payload := []byte(fmt.Sprintf("DDR|%.6f|%.6f", latitude, longitude))
	rn2483.MacTx(false, 1, payload, ddrCallback)

	// send message every 30 seconds
	timeout := time.After(time.Minute*time.Duration(timeout) + time.Second*30)
	tick := time.Tick(time.Second * 30)

	for {
		select {
		case <-timeout:
			log.Info("Timed out")
			return
		case <-tick:
			rn2483.MacTx(false, 2, []byte("a"), ddrCallback)

			//if dr := rn2483.MacGetDataRate(); dr != dataRate {
			//	log.WithField("data rate", dr).Info("new data rate")
			//	return
			//}
		}
	}
}

//type receiveCallback func(port uint8, data []byte)
func ddrCallback(port uint8, data []byte) {
	if port == 1 {
		dataString := string(data)
		if strings.HasPrefix(dataString, "DDR|") {
			dataString = strings.TrimPrefix(dataString, "DDR|")

			dr := uint8(0)

			switch dataString {
			case "7":
				dr = uint8(5)
			case "8":
				dr = uint8(4)
			case "9":
				dr = uint8(3)
			case "10":
				dr = uint8(2)
			case "11":
				dr = uint8(1)
			case "12":
				dr = uint8(0)
			}

			rn2483.MacSetDataRate(dr)

			log.WithField("data rate", dr).Info("new data rate")
		}
	}
}
