package main

import (
	"fmt"

	"github.com/andlabs/ui"

	config "github.com/timbeurskens/dynamic_gui_config"

	_ "github.com/andlabs/ui/winmanifest"
)

var ic = make(chan int)

type X struct {
	Value int
}

func (x X) Create() ui.Control {
	return ui.NewProgressBar()
}

var ButtonList = struct {
	Buttons         []func() `uiconf:"{\"vertical\":true, \"labels\":[\"Button 1\", \"Button 2 - the best\"]}"`
	Num             *int
	BoolChanIllegal chan<- bool
	IntChan         chan<- int
	Hoi             X
}{
	Buttons: []func(){
		func() {
			fmt.Println("Hello from button 1")
		},
		func() {
			fmt.Println("Howdy from button 2")
		},
		func() {
			fmt.Println("G'day from button 3")
		},
	},
	IntChan: ic,
}

func main() {
	config.Start("config", 640, 400)

	_ = config.Register("buttons", &ButtonList)

	go func() {
		for b := range ic {
			fmt.Println("int chan says: ", b)
		}
	}()

	config.Show()

	select {}
}
