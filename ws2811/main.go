package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness   = 128
	ledCounts    = 600
	ledPerStrand = 300
)

// Strand ...
type Strand struct {
	Name    string
	Reverse bool
	NumLeds int
	offset  int
}

// Render ...
func (s *Strand) Render(p Pattern) {
	// assume len p < num leds
	if s.Reverse {
		for i := 0; i < s.NumLeds; i++ {
			for _, f := range p {
				if i >= s.NumLeds {
					break
				}
				dev.Leds(0)[i+s.offset] = f
				i++
			}
		}
	} else {
		for i := 0; i < s.NumLeds; i++ {
			for _, f := range p {
				if i >= s.NumLeds {
					break
				}
				dev.Leds(0)[i+s.offset] = f
				i++
			}
		}
	}
}

// Pattern ...
type Pattern []uint32

var (
	wg      sync.WaitGroup
	dev     *ws2811.WS2811
	strands []*Strand
)

func main() {
	fmt.Println("Starting....")
	// Setup device
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	var err error
	dev, err = ws2811.MakeWS2811(&opt)
	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		panic(err)
	}
	defer dev.Fini()

	//
	// setup strands
	//
	strands = []*Strand{
		&Strand{
			Name:    "Left",
			NumLeds: 300,
			Reverse: false,
		},
		&Strand{
			Name:    "Right",
			NumLeds: 300,
			Reverse: true,
		},
	}

	//
	// Run
	//
	run()

	wg.Wait()
}

func render() {
	for _, s := range strands {
		s.Render()
	}
	dev.Render()
}

func clear() {
	for i := 0; i < ledCounts; i++ {
		dev.Leds(0)[i] = 0x000000
	}
	dev.Render()
}

func run() {
	fmt.Println("~~~ run ~~~")

	wg.Add(1)
	go func() {
		clear()

		// setRG()

		// wipe("red")
		// wipe("green")
		// wipe("white")
		wg.Done()
	}()

	// go redGreen()
	// go redWhiteGreen()
	return

	go func() {
		for i := 0; i < ledCounts; i++ {
			if i%2 == 0 {
				dev.Leds(0)[i] = 0xff0000
			} else if i%3 == 0 {
				dev.Leds(0)[i] = 0x00ff00
			} else {
				dev.Leds(0)[i] = 0x000000
			}

			dev.Render()
			time.Sleep(1 * time.Millisecond)
		}
		fmt.Println("clearing")
		for i := 0; i < ledCounts; i++ {
			dev.Leds(0)[i] = 0x000000
		}
		dev.Render()
		fmt.Println("done")
	}()

	// TODO: run async!

	// TODO: prevent running multiple at same timer

	// TODO: make service to update?

	fmt.Fprintf(w, "running...")
}

func setRG() {
	for i := 0; i < 300; i++ {
		if everyOther && i%2 != 0 {
			setC1(i, "black")
			setC2(i, "black")
			continue
		}

		if i%4 == 0 {
			setC1(i, "red")
			setC2(i, "red")
		} else {
			setC1(i, "green")
			setC2(i, "green")
		}
	}
	dev.Render()
}

func wipe(color string) {
	for i := 0; i < 300; i++ {
		if everyOther && i%2 == 0 {
			setC1(i, "black")
			setC2(i, "black")
			continue
		}

		// fmt.Println("set:", i, color)

		setC1(i, color)
		setC2(i, color)
		dev.Render()
		time.Sleep(100 * time.Millisecond)
	}
}

var everyOther = true

func setC1(i int, color string) {
	i = i + ledPerStrand
	// fmt.Println("  1:", i)
	dev.Leds(0)[i] = getHex(color)
}

func setC2(i int, color string) {
	i = ledPerStrand - i - 1
	// fmt.Println("  2:", i)
	dev.Leds(0)[i] = getHex(color)
}

func getHex(cs string) uint32 {
	switch strings.ToLower(cs) {
	case "red":
		return 0xff0000
	case "green":
		return 0x00ff00
	case "white":
		return 0xffffff
	case "black":
		return 0x000000
	}
	return 0xffffff
}

func redWhiteGreen() {
	for j := 0; j < 1000; j++ {
		for i := 0; i < 600; i++ {
			if j%3 == 0 {
				if i%3 == 0 {
					// p("red")
					dev.Leds(0)[i] = 0xff0000
				} else if i%3 == 1 {
					dev.Leds(0)[i] = 0xffffff
					// p("white")
				} else {
					dev.Leds(0)[i] = 0x00ff00
					// p("green")
				}
			} else if j%3 == 1 {
				if i%3 == 0 {
					// p("green")
					dev.Leds(0)[i] = 0x00ff00
				} else if i%3 == 1 {
					// p("red")
					dev.Leds(0)[i] = 0xff0000
				} else {
					// p("white")
					dev.Leds(0)[i] = 0xffffff
				}
			} else {
				if i%3 == 0 {
					// p("white")
					dev.Leds(0)[i] = 0xffffff
				} else if i%3 == 1 {
					// p("green")
					dev.Leds(0)[i] = 0x00ff00
				} else {
					// p("red")
					dev.Leds(0)[i] = 0xff0000
				}
			}
		}
		dev.Render()
		time.Sleep(200 * time.Millisecond)
		// fmt.Print("\r")
	}
}

func redGreen() {
	for j := 0; j < 1000; j++ {
		for i := 0; i < 600; i++ {
			if j%2 == 0 {
				if i%2 == 0 {
					dev.Leds(0)[i] = 0xff0000
				} else {
					dev.Leds(0)[i] = 0x00ff00
				}
			} else {
				if i%2 == 0 {
					dev.Leds(0)[i] = 0x00ff00
				} else {
					dev.Leds(0)[i] = 0xff0000
				}
			}
		}
		dev.Render()
		time.Sleep(200 * time.Millisecond)
	}
}
