package main

import (
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
	SomeFloatValue float64 `uiconf:"{\"name\":\"float value:\",\"min\":0,\"max\":1,\"resolution\":100}"`
	SomeBool       bool
	SomeInt        IntNotifier
	SomeString     string `uiconf:"{\"name\":\"enter text:\"}"`
}

var Config1 = ApplicationConfig{
	SomeFloatValue: 0.5,
	SomeBool:       true,
	SomeInt:        0,
	SomeString:     "default value",
}

var Config2 = ApplicationConfig{}

func main() {
	config.Start("config", 640, 400)

	_ = config.Register("panel1", &Config1)
	_ = config.Register("panel2", &Config2)

	config.Show()

	ticker := time.Tick(time.Second)

	for range ticker {
		log.Println("1", Config1)
		log.Println("2", Config2)
	}
}
