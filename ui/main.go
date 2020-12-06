package main

import (
	"fmt"

	"github.com/jaredwarren/app"
	"github.com/jaredwarren/blink/service"
)

func main() {
	conf := &app.WebConfig{
		Host: "127.0.0.1",
		Port: 8084,
	}
	a := app.NewWeb(conf)

	service.Register(a, nil)

	d := <-a.Exit
	fmt.Printf("Done:%+v\n", d)
}
