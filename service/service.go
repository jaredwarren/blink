package service

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jaredwarren/app"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 128
	ledCounts  = 600
)

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
	tmpl := template.Must(template.ParseFiles("layout.html", "base.html"))
	tmpl.ExecuteTemplate(w, "base", struct {
		PageTitle string
	}{
		"Something",
	})
}

func (s *Server) clear(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < ledCounts; i++ {
		s.dev.Leds(0)[i] = 0x000000
	}
	s.dev.Render()
}

func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running...")
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
		for i := 0; i < ledCounts; i++ {
			s.dev.Leds(0)[i] = 0x000000
		}
		s.dev.Render()
		fmt.Println("done")
	}()
	fmt.Fprintf(w, "Running...")
}
