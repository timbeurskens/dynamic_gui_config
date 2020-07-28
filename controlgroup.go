package dynamic_gui_config

import "github.com/andlabs/ui"

type controlGroup []ValueControl

func (c controlGroup) Create() ui.Control {
	container := ui.NewVerticalBox()
	container.SetPadded(true)

	for _, field := range c {
		container.Append(field.Create(), false)
	}

	return container
}
