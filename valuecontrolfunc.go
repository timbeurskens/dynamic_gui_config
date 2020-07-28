package dynamic_gui_config

import "github.com/andlabs/ui"

type ValueControlFunc func() ui.Control

func (v ValueControlFunc) Create() ui.Control {
	return v()
}
