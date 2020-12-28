# Dynamic Config for Go

This go package gives users the option to dynamically adjust variables in an graphical environment.
It uses the [andlabs ui library](https://github.com/andlabs/ui) for showing the GUI.

**Supported platforms**

The package OS support is restricted only by its UI framework. In the future it might be possible to create an platform agnostic library such that multiple UI libraries can be used.
For now only Windows and Linux on x86-64 are supported.
On Windows, make sure cgo is supported.

## Installation

Install the package by running:

```bash
go get github.com/timbeurskens/dynamic_gui_config
```

## Usage

Examples can be found in the `examples` directory.

```go
type ApplicationConfig struct {
	SomeFloatValue    float64 `uiconf:"{\"name\":\"float value:\",\"min\":0,\"max\":1,\"resolution\":100}"`
	SomeBool          bool
}

var Config = ApplicationConfig{
	SomeFloatValue: 0.5,
	SomeBool:       true,
}

func main() {
	config.Start("config", 640, 400)

	_ = config.Register("config", &Config)

	config.Show()

	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		log.Println(Config)
	}
}
```

## Supported types

### Custom types



## Known issues

### Exit code -1073741511 (0xC0000139) on Windows

**Solution:** make sure to include a windows manifest file your application (e.g. `import _ github.com/andlabs/ui/winmanifest`)