package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jaredwarren/app"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cobra"
)

const (
	brightness   = 128
	ledCounts    = 600
	ledPerStrand = 300
)

// Strand ...
type Strand struct {
	Reverse bool
	NumLeds int
}

// Frame ...
type Frame struct {
	Pattern []uint32
}

// // Device should match https://pkg.go.dev/github.com/rpi-ws281x/rpi-ws281x-go#MakeWS2811
// type Device interface {
// 	Fini()
// 	Init() error
// 	Leds(channel int) []uint32
// 	Render() error
// 	SetBrightness(channel int, brightness int)
// 	SetLedsSync(channel int, leds []uint32) error
// 	Wait() error
// }

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var wg sync.WaitGroup

func main() {
	fmt.Println("Starting....")
	// Setup device
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		panic(err)
	}
	defer dev.Fini()

	// remove after testing
	s := Server{dev}
	s.run(nil, nil)

	wg.Wait()

	return

	// setup server
	conf := &app.WebConfig{
		Host: ":8084",
	}
	a := app.NewWeb(conf)

	Register(a, dev)
	fmt.Println("Ready....")
	d := <-a.Exit
	fmt.Printf("Done:%+v\n", d)

	return
	// old cmd stuff

	rootCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear Everything",
		Long:  `clear`,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < ledCounts; i++ {
				dev.Leds(0)[i] = 0x000000
			}
			checkError(dev.Render())
		},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "show",
		Long:  `show`,
		Run: func(cmd *cobra.Command, args []string) {
			for {
				chase(dev)
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
			clear(dev)
		},
	})
	rootCmd.Execute()
}

func chase(dev *ws2811.WS2811) {
	for i := 0; i < ledCounts; i++ {
		if i%2 == 0 {
			dev.Leds(0)[i] = 0xff007f
		} else if i%3 == 0 {
			dev.Leds(0)[i] = 0xff007f
		} else {
			dev.Leds(0)[i] = 0xff007f
		}

		checkError(dev.Render())
		time.Sleep(1 * time.Millisecond)
	}
	for i := 0; i < ledCounts; i++ {
		dev.Leds(0)[i] = 0x000000
	}
	checkError(dev.Render())
}

func clear(dev *ws2811.WS2811) {
	for i := 0; i < ledCounts; i++ {
		dev.Leds(0)[i] = 0x000000
	}
	checkError(dev.Render())
}

func altrg(dev *ws2811.WS2811) {
	for i := 0; i < ledCounts; i++ {
		if i%2 == 0 {
			dev.Leds(0)[i] = 0xff0000
		} else {
			dev.Leds(0)[i] = 0x000000
		}
	}
	checkError(dev.Render())
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < ledCounts; i++ {
		if i%2 == 0 {
			dev.Leds(0)[i] = 0x000000
		} else {
			dev.Leds(0)[i] = 0x00ff00
		}
	}
	checkError(dev.Render())
	time.Sleep(200 * time.Millisecond)
}

// Server ..
type Server struct {
	dev *ws2811.WS2811
}

// Register ...
func Register(a *app.Service, dev *ws2811.WS2811) {
	m := a.Mux

	s := Server{dev}

	m.HandleFunc("/", home).Methods("GET")
	m.HandleFunc("/clear", s.clear).Methods("GET")
	m.HandleFunc("/run", s.run).Methods("GET")

}

// Close handler.
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("~~~ Home ~~~")
	tmpl := template.Must(template.ParseFiles("layout.html", "base.html"))
	tmpl.ExecuteTemplate(w, "base", struct {
		PageTitle string
	}{
		"Something",
	})
}

func (s *Server) clear(w http.ResponseWriter, r *http.Request) {
	fmt.Println("~~~ clear ~~~")
	for i := 0; i < ledCounts; i++ {
		s.dev.Leds(0)[i] = 0x000000
	}
	s.dev.Render()
	fmt.Fprintf(w, "clear")
}

func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	fmt.Println("~~~ run ~~~")

	wg.Add(1)
	go func() {
		clear(s.dev)

		// s.setRG()

		// s.wipe("red")
		// s.wipe("green")
		// s.wipe("white")
		wg.Done()
	}()

	// go s.redGreen()
	// go s.redWhiteGreen()
	return

	go func() {
		for i := 0; i < ledCounts; i++ {
			if i%2 == 0 {
				s.dev.Leds(0)[i] = 0xff0000
			} else if i%3 == 0 {
				s.dev.Leds(0)[i] = 0x00ff00
			} else {
				s.dev.Leds(0)[i] = 0x000000
			}

			s.dev.Render()
			time.Sleep(1 * time.Millisecond)
		}
		fmt.Println("clearing")
		for i := 0; i < ledCounts; i++ {
			s.dev.Leds(0)[i] = 0x000000
		}
		s.dev.Render()
		fmt.Println("done")
	}()

	// TODO: run async!

	// TODO: prevent running multiple at same timer

	// TODO: make service to update?

	fmt.Fprintf(w, "running...")
}

func (s *Server) setRG() {
	for i := 0; i < 300; i++ {
		if everyOther && i%2 != 0 {
			s.setC1(i, "black")
			s.setC2(i, "black")
			continue
		}

		if i%4 == 0 {
			s.setC1(i, "red")
			s.setC2(i, "red")
		} else {
			s.setC1(i, "green")
			s.setC2(i, "green")
		}
	}
	s.dev.Render()
}

func (s *Server) wipe(color string) {
	for i := 0; i < 300; i++ {
		if everyOther && i%2 == 0 {
			s.setC1(i, "black")
			s.setC2(i, "black")
			continue
		}

		// fmt.Println("set:", i, color)

		s.setC1(i, color)
		s.setC2(i, color)
		s.dev.Render()
		time.Sleep(100 * time.Millisecond)
	}
}

var everyOther = true

func (s *Server) setC1(i int, color string) {
	i = i + ledPerStrand
	// fmt.Println("  1:", i)
	s.dev.Leds(0)[i] = getHex(color)
}

func (s *Server) setC2(i int, color string) {
	i = ledPerStrand - i - 1
	// fmt.Println("  2:", i)
	s.dev.Leds(0)[i] = getHex(color)
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

func (s *Server) redWhiteGreen() {
	for j := 0; j < 1000; j++ {
		for i := 0; i < 600; i++ {
			if j%3 == 0 {
				if i%3 == 0 {
					// p("red")
					s.dev.Leds(0)[i] = 0xff0000
				} else if i%3 == 1 {
					s.dev.Leds(0)[i] = 0xffffff
					// p("white")
				} else {
					s.dev.Leds(0)[i] = 0x00ff00
					// p("green")
				}
			} else if j%3 == 1 {
				if i%3 == 0 {
					// p("green")
					s.dev.Leds(0)[i] = 0x00ff00
				} else if i%3 == 1 {
					// p("red")
					s.dev.Leds(0)[i] = 0xff0000
				} else {
					// p("white")
					s.dev.Leds(0)[i] = 0xffffff
				}
			} else {
				if i%3 == 0 {
					// p("white")
					s.dev.Leds(0)[i] = 0xffffff
				} else if i%3 == 1 {
					// p("green")
					s.dev.Leds(0)[i] = 0x00ff00
				} else {
					// p("red")
					s.dev.Leds(0)[i] = 0xff0000
				}
			}
		}
		s.dev.Render()
		time.Sleep(200 * time.Millisecond)
		// fmt.Print("\r")
	}
}

func (s *Server) redGreen() {
	for j := 0; j < 1000; j++ {
		for i := 0; i < 600; i++ {
			if j%2 == 0 {
				if i%2 == 0 {
					s.dev.Leds(0)[i] = 0xff0000
				} else {
					s.dev.Leds(0)[i] = 0x00ff00
				}
			} else {
				if i%2 == 0 {
					s.dev.Leds(0)[i] = 0x00ff00
				} else {
					s.dev.Leds(0)[i] = 0xff0000
				}
			}
		}
		s.dev.Render()
		time.Sleep(200 * time.Millisecond)
	}
}
