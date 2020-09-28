package dynamic_gui_config

import "github.com/andlabs/ui"

type LabeledGuiField struct {
	Label   ValueControl
	Factory ValueControl
}

func (l LabeledGuiField) Create() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	// Label should not stretch
	hbox.Append(l.Label.Create(), false)

	// Factory should fill available space
	hbox.Append(l.Factory.Create(), true)

	return hbox
}
