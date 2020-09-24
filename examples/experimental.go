package main

import (
	"fmt"

	config "github.com/timbeurskens/dynamic_gui_config"

	_ "github.com/andlabs/ui/winmanifest"
)

var ButtonList = struct {
	Buttons []func() `uiconf:"{\"vertical\":true}"`
	Num     *int
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
}

func main() {
	config.Start("config", 640, 400)

	_ = config.Register("buttons", &ButtonList)

	config.Show()

	select {}
}
