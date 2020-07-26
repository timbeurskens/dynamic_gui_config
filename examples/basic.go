package main

import (
	"fmt"
	"log"
	"time"

	config "github.com/timbeurskens/dynamic_gui_config"

	_ "github.com/andlabs/ui/winmanifest"
)

type IntNotifier int

func (i IntNotifier) OnValueChanged() {
	log.Println("value changed: ", i)
}

type ApplicationConfig struct {
	SomeFloatValue    float64 `uiconf:"{\"name\":\"float value:\",\"min\":0,\"max\":1,\"resolution\":100}"`
	SomeBool          bool
	SomeInt           IntNotifier
	SomeString        string `uiconf:"{\"name\":\"enter text:\"}"`
	SomeFunction      func()
	SomeOtherFunction func(int) int
	SomeUnsignedInt   uint `uiconf:"{\"max\":200}"`
	notExported       int
}

var Config1 = ApplicationConfig{
	SomeFloatValue: 0.5,
	SomeBool:       true,
	SomeInt:        0,
	SomeString:     "default value",
}

var Config2 = ApplicationConfig{
	SomeFunction: func() {
		fmt.Println("Hello, world!")
	},
}

func main() {
	config.Start("config", 640, 400)

	_ = config.Register("panel1", &Config1)
	_ = config.Register("panel2", &Config2)

	config.Show()

	ticker := time.NewTicker(time.Second)

	go func() {
		<-time.After(10 * time.Second)
		ticker.Stop()
	}()

	for range ticker.C {
		log.Println("1", Config1)
		log.Println("2", Config2)
	}

	config.Stop()
}
