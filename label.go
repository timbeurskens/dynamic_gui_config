package dynamic_gui_config

import "github.com/andlabs/ui"

type Label string

func (l Label) Create() ui.Control {
	return ui.NewLabel(string(l))
}
