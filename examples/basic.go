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
	FloatArray        [3]float64
	SomeSubStruct     struct {
		SomeSubInt int
		AnotherInt int
		SomeBool   bool
	}
}

var Config1 = ApplicationConfig{
	SomeFloatValue: 0.5,
	SomeBool:       true,
	SomeInt:        0,
	SomeString:     "default value",
	SomeFunction: func() {
		fmt.Println("Hello, world!")
	},
}

func main() {
	config.Start("config", 640, 400)

	_ = config.Register("panel1", &Config1)

	config.Show()

	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		log.Println(Config1)
	}
}
