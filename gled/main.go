package main

import (
	"encoding/json"
	"fmt"
	"gocron-master/gocron-master"
	"html/template"
	"log"
	"net/http"
	"os"
	"rpi_ws281x/golang/ws2811"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var buildStable bool = true
var alerting bool = false
var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	pin = rpio.Pin(10)
)
var buildch chan bool

func main() {

	if err := ws2811.Init(18, 300, 32); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buildch := make(chan bool)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	fmt.Println("Initialized")

	rainbowCycle()

	http.HandleFunc("/", Home)

	http.HandleFunc("/alertBuildUnstable", func(w http.ResponseWriter, r *http.Request) {

		handleBuildError(w, buildch)
	})

	http.HandleFunc("/alertTestUnstable", func(w http.ResponseWriter, r *http.Request) {

		handleTestError(w, buildch)
	})

	http.HandleFunc("/alertBuildStable", func(w http.ResponseWriter, r *http.Request) {

		handleBuildStable(w, buildch)
	})

	http.HandleFunc("/alertStop", func(w http.ResponseWriter, r *http.Request) {

		handleStop(w, buildch)
	})

	go scheduleJobs(buildch)
	log.Fatal(http.ListenAndServe(":8083", nil))
	// remove, clear and next_run

}

func scheduleJobs(buildch chan bool) {

	//schedule some jobs
	// gocron.Every(1).Day().At("18:28").Do(standUp)
	gocron.Every(1).Day().At("10:30").Do(standUp, buildch)
	gocron.Every(1).Day().At("11:45").Do(standUp, buildch)
	gocron.Every(1).Day().At("14:00").Do(standUp, buildch)
	// gocron.Every(1).Monday().At("18:30").Do(task)
	_, time := gocron.NextRun()
	fmt.Println("next run: ")
	fmt.Println(time)
	// function Start start all the pending jobs
	<-gocron.Start()
}

func standUp(buildch chan bool) {
	if alerting == true {
		alerting = false
		buildch <- true
	}

	stripStandUp(getColour(255, 255, 255))
	ws2811.Clear()

	if buildStable == true {
		go setStripColour(254, 0, 0)
	} else {
		alerting = true
		go toggleTestAlert(buildch)
	}

}

func Home(w http.ResponseWriter, req *http.Request) {
	render(w, "/home/pi/gowork/alert/alert.html")
}

func render(w http.ResponseWriter, tmpl string) {
	tmpl = fmt.Sprintf("%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func handleBuildError(w http.ResponseWriter, buildch chan bool) {

	mapResponse, _ := json.Marshal("unstable")

	fmt.Fprintf(w, "%q", string(mapResponse))

	time.Sleep(time.Second)

	buildStable = false
	if alerting == false {
		alerting = true
		go toggleBuildAlert(buildch)
	}

}

func handleTestError(w http.ResponseWriter, buildch chan bool) {

	mapResponse, _ := json.Marshal("unstableTest")

	fmt.Fprintf(w, "%q", string(mapResponse))

	time.Sleep(time.Second)

	buildStable = false
	if alerting == false {
		alerting = true
		go toggleTestAlert(buildch)
	}

}

func toggleBuildAlert(buildch chan bool) {
	buildStable = false
	ws2811.Clear()
	for buildStable == false {

		stripStart(getColour(0, 255, 0))

		//  pin.Toggle()
		//              time.Sleep(time.Second / 6)
		// use select for non blocking operation
		select {
		case buildStable = <-buildch:
			fmt.Println("build stable? : ", buildStable)
		default:
			fmt.Println("build unstable!")
		}
	}
}

func toggleTestAlert(buildch chan bool) {
	ws2811.Clear()
	buildStable = false
	for buildStable == false {
		stripStart(getColour(165, 255, 0))

		//  pin.Toggle()
		//                time.Sleep(time.Second)
		// use select for non blocking operation
		select {
		case buildStable = <-buildch:
			fmt.Println("build stable? : ", buildStable)
		default:
			fmt.Println("build unstable!")
		}
	}
}

func handleBuildStable(w http.ResponseWriter, buildch chan bool) {

	mapResponse, _ := json.Marshal("stable")
	time.Sleep(time.Second)
	fmt.Fprintf(w, "%q", string(mapResponse))
	alerting = false
	if buildStable == false {
		buildch <- true
		setStripColour(254, 0, 0)
		buildStable = true
	} else {
		//be proud make it blue
		setStripColour(0, 0, 254)
	}
}

func handleStop(w http.ResponseWriter, buildch chan bool) {

	mapResponse, _ := json.Marshal("stop")
	time.Sleep(time.Second)
	fmt.Fprintf(w, "%q", string(mapResponse))
	if buildStable == false {
		buildch <- true
	}
	time.Sleep(time.Second)
	ws2811.Clear()
	ws2811.Render()
	buildStable = true
	alerting = false
}

func stripStart(colour uint32) {
	for i := 0; i < 150; i++ {
		ws2811.SetLed(i, colour)
		ws2811.SetLed(299-i, colour)
		ws2811.Render()
		time.Sleep(time.Second / 60)
	}

	for i := 149; i >= 0; i-- {
		ws2811.SetLed(i, 0)
		ws2811.SetLed(149+(149-i), 0)
		ws2811.Render()
		time.Sleep(time.Second / 60)
	}

}

func stripStandUp(colour uint32) {
	for j := 0; j < 20; j++ {
		for i := 0; i < 300; i++ {
			ws2811.SetLed(i, colour)
		}
		ws2811.Render()
		time.Sleep(time.Second / 2)
		for i := 0; i < 300; i++ {
			ws2811.SetLed(i, 0)
		}
		ws2811.Render()
		time.Sleep(time.Second / 2)
	}
}

func wheel(pos int) uint32 {
	//Generate rainbow colors across 0-255 positions
	var rgb uint32

	if pos < 85 {
		rgb = uint32(pos) * 3
		rgb = (rgb << 8) + 255 - uint32(pos)
		rgb = (rgb << 8) + 0
		return rgb
	} else if pos < 170 {
		pos -= 85
		rgb = 255 - uint32(pos)*3
		rgb = (rgb << 8) + 0
		rgb = (rgb << 8) + uint32(pos)*3
		return rgb
	} else {
		pos -= 170
		rgb = 0
		rgb = (rgb << 8) + uint32(pos)*3
		rgb = (rgb << 8) + 255 - uint32(pos)*3
		return rgb
	}
}

func getColour(green uint32, red uint32, blue uint32) uint32 {
	var rgb uint32
	rgb = green
	rgb = (rgb << 8) + red
	rgb = (rgb << 8) + blue
	return rgb
}

func setStripColour(green uint32, red uint32, blue uint32) {
	for i := 0; i < 300; i++ {
		ws2811.SetLed(i, getColour(green, red, blue))
	}
	ws2811.Render()
}

func rainbowCycle() {
	//Draw rainbow that uniformly distributes itself across all pixels.
	var j int = 0
	//for j := 0; j <= 256 * 1 ; j++ {

	for i := 0; i < 300; i++ {
		ws2811.SetLed(i, wheel((int(i*256/300)+j)&255))
		ws2811.Render()
		time.Sleep(time.Second / 100)
	}
	//}
	setStripColour(0, 0, 254)
}
