package dynamic_gui_config

import "github.com/andlabs/ui"

type structGuiField struct {
	Properties StructTagProperties
	Factory    ValueControl
}

func (s structGuiField) Create() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	hbox.Append(s.Properties.Name.Create(), false)
	hbox.Append(s.Factory.Create(), true)

	return hbox
}
