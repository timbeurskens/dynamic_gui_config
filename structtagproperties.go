package dynamic_gui_config

import (
	"encoding/json"
	"errors"
)

const (
	tagkey = "uiconf"
)

// StructTagProperties represents the values supported in a uiconf struct tag
type StructTagProperties struct {
	Name       Label  `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Min        int    `json:"min,omitempty"`
	Max        int    `json:"max,omitempty"`
	Resolution int    `json:"resolution,omitempty"`
	Horizontal bool   `json:"horizontal,omitempty"`
}

var structTagDefaults = StructTagProperties{
	Name:       "",
	Type:       "",
	Min:        0,
	Max:        100,
	Resolution: 10,
	Horizontal: false,
}

func DefaultProperties() StructTagProperties {
	return structTagDefaults
}

// Parse decodes the json format tag string into the current properties pointer
func (properties *StructTagProperties) Parse(tag string) error {
	if tag == "" {
		return errors.New("empty tag, using defaults")
	}

	return json.Unmarshal([]byte(tag), properties)
}

// ParseStructTag returns a new StructTagProperties object with the contents from tag
func ParseStructTag(tag string) (properties StructTagProperties, err error) {
	properties = structTagDefaults

	err = properties.Parse(tag)
	return
}
