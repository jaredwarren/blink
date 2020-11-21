package main

import (
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cobra"
)

const (
	brightness = 128
	ledCounts  = 300
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
	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		panic(err)
	}
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
}
