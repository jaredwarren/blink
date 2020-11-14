// Copyright 2018 Jacques Supcik / HEIA-FR
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cobra"
)

const (
	brightness = 128
	width      = 8
	height     = 8
	ledCounts  = width * height
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	opt := ws2811.DefaultOptions
	// opt.Channels[0].Brightness = 0
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = 300

	dev, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	checkError(dev.Init())
	defer dev.Fini()

	rootCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear Everything",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < 300; i++ {
				dev.Leds(0)[i] = 0x000000
			}
			checkError(dev.Render())
		},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Clear Everything",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < 300; i++ {
				if i%2 == 0 {
					dev.Leds(0)[i] = 0xff0000
				} else if i%3 == 0 {
					dev.Leds(0)[i] = 0x00ff00
				} else {
					dev.Leds(0)[i] = 0xffffff
				}

				checkError(dev.Render())
				time.Sleep(100 * time.Millisecond)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "clear",
		Short: "Clear Everything",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < 300; i++ {
				dev.Leds(0)[i] = 0x000000
			}
			checkError(dev.Render())
		},
	})
	rootCmd.Execute()

	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		color := uint32(0xff0000)
	// 		if x > 2 && x < 5 && y > 0 && y < 7 {
	// 			color = 0xffffff
	// 		}
	// 		if x > 0 && x < 7 && y > 2 && y < 5 {
	// 			color = 0xffffff
	// 		}
	// 		dev.Leds(0)[x*height+y] = color
	// 	}
	// }
	// checkError(dev.Render())
}
