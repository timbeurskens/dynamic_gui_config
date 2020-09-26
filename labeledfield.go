package dynamic_gui_config

import "github.com/andlabs/ui"

type LabeledGuiField struct {
	Label   Label
	Factory ValueControl
}

func (l LabeledGuiField) Create() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	hbox.Append(l.Label.Create(), false)
	hbox.Append(l.Factory.Create(), true)

	return hbox
}
