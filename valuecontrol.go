package dynamic_gui_config

import "github.com/andlabs/ui"

// ValueControl specifies an object which can create a ui.Control
type ValueControl interface {
	Create() ui.Control
}
