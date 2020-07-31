package dynamic_gui_config

import "github.com/andlabs/ui"

// ValueControlFunc accepts functions that create ui.Control objects
// the wrapping functions will be run in the ui thread
type ValueControlFunc func() ui.Control

func (v ValueControlFunc) Create() ui.Control {
	return v()
}
