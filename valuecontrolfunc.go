package dynamic_gui_config

import "github.com/andlabs/ui"

// ValueControlFunc accepts functions that create ui.Control objects
// the wrapping functions will be run in the ui thread
type ValueControlFunc func() ui.Control

// Create implements ValueControl
func (v ValueControlFunc) Create() ui.Control {
	return v()
}

// StrictValueControlFunc is a variant of ValueControlFunc which panics when the function returns an error
type StrictValueControlFunc func() (ui.Control, error)

// Create implements ValueControl
func (s StrictValueControlFunc) Create() ui.Control {
	if ctrl, err := s(); err != nil {
		panic(err)
	} else {
		return ctrl
	}
}
