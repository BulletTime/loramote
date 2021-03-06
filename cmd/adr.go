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
	"time"

	"github.com/apex/log"
	"github.com/bullettime/rn2483"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	timeout   uint8
	dataRate  uint8
	confirmed bool
)

// adrCmd represents the adr command
var adrCmd = &cobra.Command{
	Use:   "adr",
	Short: "Test ADR algorithm",
	Run: func(cmd *cobra.Command, args []string) {
		TestADR()
	},
}

func init() {
	RootCmd.AddCommand(adrCmd)

	adrCmd.Flags().Uint8VarP(&dataRate, "datarate", "r", 0, "set data rate (default 0)")
	adrCmd.Flags().Uint8VarP(&timeout, "timeout", "t", 10, "set timeout in minutes")
	adrCmd.Flags().BoolVarP(&confirmed, "confirmed", "c", false, "use confirmed uplink")
}

func TestADR() {
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
	rn2483.MacSetADR(true)
	//linkchk := uint16(660)
	//rn2483.MacSetLinkCheck(linkchk)

	log.WithFields(log.Fields{
		"data rate": rn2483.MacGetDataRate(),
		"power":     rn2483.MacGetPowerIndex(),
		"adr":       rn2483.MacGetADR(),
		//"link check": linkchk,
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

	// send message every 30 seconds
	timeout := time.After(time.Minute*time.Duration(timeout) + time.Second*30)
	tick := time.Tick(time.Second * 30)

	counter := 0

	for {
		select {
		case <-timeout:
			log.Info("Timed out")
			return
		case <-tick:
			rn2483.MacTx(confirmed, 2, []byte("a"), nil)
			log.WithFields(log.Fields{
				"confirmed": confirmed,
				"port":      2,
				"data":      "a",
				"frame":     counter,
			}).Info("uplink")

			counter++
		}
	}
}

type myLogger struct {
	level string
}

func (l myLogger) Println(v ...interface{}) {
	switch l.level {
	case "ERROR":
		log.Error(fmt.Sprintln(v...))
	case "WARN":
		log.Warn(fmt.Sprintln(v...))
	case "INFO":
		log.Info(fmt.Sprintln(v...))
	case "DEBUG":
		log.Debug(fmt.Sprintln(v...))
	}
}

func (l myLogger) Printf(format string, v ...interface{}) {
	switch l.level {
	case "ERROR":
		log.Errorf(format, v...)
	case "WARN":
		log.Warnf(format, v...)
	case "INFO":
		log.Infof(format, v...)
	case "DEBUG":
		log.Debugf(format, v...)
	}
}
