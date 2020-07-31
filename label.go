package dynamic_gui_config

import "github.com/andlabs/ui"

// Label represents a ui.Label type.
type Label string

// Create creates the ui.Control object in the ui thread
func (l Label) Create() ui.Control {
	return ui.NewLabel(string(l))
}
