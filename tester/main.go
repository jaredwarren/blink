package main

import (
	"fmt"
	"math"
	"time"

	. "github.com/logrusorgru/aurora"
)

const (
	// NumLed ...
	NumLed = 150
)

// Printer ...
type Printer interface {
	Print(Frame)
}

// Strand ...
type Strand struct {
	Name    string
	Reverse bool
	NumLeds int
	offset  int
	Printer
}

// Render ...
func (s *Strand) Render(f Frame) {

	if s.Reverse {
		for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
			f[i], f[j] = f[j], f[i]
		}
	}
	s.Print(f)

	// // assume len p < num leds
	// out := make([]string, s.NumLeds)

	// if s.Reverse {
	// 	for i := 0; i < s.NumLeds; i++ {

	// 		// for _, f := range p {
	// 		// 	if i >= s.NumLeds {
	// 		// 		break
	// 		// 	}
	// 		// 	out[i] = f
	// 		// 	// dev.Leds(0)[i+s.offset] = f
	// 		// 	i++
	// 		// }
	// 	}
	// } else {
	// 	for i := s.NumLeds - 1; i >= 0; i-- {
	// 		for _, f := range p {
	// 			if i < 0 {
	// 				break
	// 			}
	// 			// dev.Leds(0)[i+s.offset] = f
	// 			out[i] = f
	// 			i--
	// 		}
	// 	}
	// }

	// for _, c := range out {
	// 	pp(c)
	// }
}

// Pattern ...
// type Pattern []uint32
type Pattern []string

// Frame ... len == num leds
type Frame []string

// Push ...
func (f *Frame) Push(c string) bool {
	if len(*f) == cap(*f) {
		return false
	}
	*f = append(*f, c)
	return true
}

// PushP ...
func (f *Frame) PushP(p Pattern) {
	for _, pp := range p {
		if ok := f.Push(pp); !ok {
			return
		}
	}
}

// Fill ...
func (f *Frame) Fill(p Pattern) {
	// push until full
	for len(*f) < cap(*f) {
		for _, pp := range p {
			if ok := f.Push(pp); !ok {
				return
			}
		}
	}
}

// Print ...
func (f *Frame) Print() {
	for _, p := range *f {
		pp(p)
	}

}

// NewFrame ...
func NewFrame(l int) *Frame {
	f := make(Frame, 0, l)
	return &f
}

func main() {

	// ff := make(Frame, 0, 10)
	// ff.PushP(Pattern{"a", "b", "c"})
	// ff.PushP(Pattern{"a", "b", "c"})
	// ff.PushP(Pattern{"a", "b", "c"})
	// ff.PushP(Pattern{"a", "b", "c"})
	// ff.PushP(Pattern{"a", "b", "c"})
	// ff.Fill(Pattern{"a", "b", "c"})

	// for i := 0; i < 14; i++ {
	// 	fmt.Println(i)
	// 	ok := ff.Push(fmt.Sprintf("%d", i))
	// 	fmt.Println(" ", ok)
	// }

	// fmt.Printf("%+v\n", ff)

	// return

	// fmt.Println("Hello,", Magenta("Aurora"))
	// fmt.Println(Bold(Red(".")))
	fs := genWipe(10, Pattern{"red", "white"})

	for _, f := range fs {
		for _, ff := range f {
			pp(ff)
		}

		fmt.Println("")
	}
	// fmt.Println(fs)

	fmt.Println()

	{
		f := NewFrame(100)
		f.Fill(Pattern{"red", "green", "white"})
		f.Print()
	}
	fmt.Println()

	l := &Strand{
		Name:    "Left",
		NumLeds: 300,
		Reverse: false,
		Printer: &ConsolePrinter{},
	}
	r := &Strand{
		Name:    "Right",
		NumLeds: 300,
		Reverse: true,
		Printer: &ConsolePrinter{},
	}
	fmt.Println("\n---------")
	{
		f := NewFrame(100)
		f.Fill(Pattern{"red", "green", "white"})
		l.Render(*f)
	}
	fmt.Println("\n---------")
	{
		f := NewFrame(100)
		f.Fill(Pattern{"red", "green", "white"})
		r.Render(*f)
	}
	fmt.Println("\n---------")

	// c1 := []string{"red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green", "red", "white", "green"}
	// c2 := []string{"red", "red", "green", "white"}

	// print(c1, c2)

}

// It would be nice to  "spread out" / repeat frames
func genWipe(l int, p Pattern) []Frame {
	frames := make([]Frame, l)
	for i := 0; i < l; i++ {
		f := make(Frame, 0, l)

		done := true // there's got to be a bette way
		max := int(math.Ceil(float64(l) / float64(len(p))))
		for j := 0; j < max; j++ {
			if j <= i {
				f.PushP(p)
			} else {
				f.Fill(Pattern{"green"}) // TODO: replace with "nil"
				done = false
			}
		}
		frames[i] = f
		if done {
			break
		}
	}

	return frames
}

//
//

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
		pp(s)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n-----------")

	fmt.Println("\n-----------")

	for _, s := range c2 {
		pp(s)
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("")

}

// ConsolePrinter console printer
type ConsolePrinter struct{}

// Print ...
func (p *ConsolePrinter) Print(f Frame) {
	for _, p := range f {
		pp(p)
	}
}

func pp(x string) {
	switch x {
	case "red":
		fmt.Print(Bold(Red("⬤")))
	case "green":
		fmt.Print(Bold(Green("⬤")))
	case "black":
		fmt.Print(Bold(Black("⬤")))
	case "off":
		fmt.Print(Bold(Black("⬤")))
	case "white":
		fmt.Print(Bold("⬤"))
	default:
		fmt.Print(Bold(" "))
	}
}

func wipe(color string) {
	fmt.Print("\r")
	for i := 0; i < 155; i++ {
		pp(color)
		time.Sleep(50 * time.Millisecond)
	}
}

func redGreen() {
	fmt.Print("\r")
	for j := 0; j < 10; j++ {
		for i := 0; i < 155; i++ {
			if j%2 == 0 {
				if i%2 == 0 {
					pp("red")
				} else {
					pp("green")
				}
			} else {
				if i%2 == 0 {
					pp("green")
				} else {
					pp("red")
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
					pp("red")
				} else if i%3 == 1 {
					pp("white")
				} else {
					pp("green")
				}
			} else if j%3 == 1 {
				if i%3 == 0 {
					pp("green")
				} else if i%3 == 1 {
					pp("red")
				} else {
					pp("white")
				}
			} else {
				if i%3 == 0 {
					pp("white")
				} else if i%3 == 1 {
					pp("green")
				} else {
					pp("red")
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
		fmt.Print("\r")
	}
}
