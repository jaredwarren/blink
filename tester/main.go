package main

import (
	"fmt"
	"time"

	. "github.com/logrusorgru/aurora"
)

const (
	// NumLed ...
	NumLed = 150
)

func main() {
	// fmt.Println("Hello,", Magenta("Aurora"))
	// fmt.Println(Bold(Red(".")))

	c1 := []string{"red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green"}
	c2 := []string{"red", "red", "green", "white"}

	print(c1, c2)

}

func print(c1, c2 []string) {
	// reverse c2
	for i, j := 0, len(c2)-1; i < j; i, j = i+1, j-1 {
		c2[i], c2[j] = c2[j], c2[i]
	}

	redWhiteGreen()
	wipe("green")
	redGreen()
	wipe("red")
	wipe("white")
	redGreen()

	fmt.Println("")
	for _, s := range c1 {
		p(s)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n-----------")

	fmt.Println("\n-----------")

	for _, s := range c2 {
		p(s)
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("")

}

func p(x string) {
	switch x {
	case "red":
		fmt.Print(Bold(Red(".")))
	case "green":
		fmt.Print(Bold(Green(".")))
	case "white":
		fmt.Print(Bold("."))

	}
}

func wipe(color string) {
	fmt.Print("\r")
	for i := 0; i < 155; i++ {
		p(color)
		time.Sleep(50 * time.Millisecond)
	}
}

func redGreen() {
	fmt.Print("\r")
	for j := 0; j < 10; j++ {
		for i := 0; i < 155; i++ {
			if j%2 == 0 {
				if i%2 == 0 {
					p("red")
				} else {
					p("green")
				}
			} else {
				if i%2 == 0 {
					p("green")
				} else {
					p("red")
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
		fmt.Print("\r")
	}
}

func redWhiteGreen() {
	fmt.Print("\r")
	for j := 0; j < 10; j++ {
		for i := 0; i < 155; i++ {
			if j%3 == 0 {
				if i%3 == 0 {
					p("red")
				} else if i%3 == 1 {
					p("white")
				} else {
					p("green")
				}
			} else if j%3 == 1 {
				if i%3 == 0 {
					p("green")
				} else if i%3 == 1 {
					p("red")
				} else {
					p("white")
				}
			} else {
				if i%3 == 0 {
					p("white")
				} else if i%3 == 1 {
					p("green")
				} else {
					p("red")
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
		fmt.Print("\r")
	}
}
