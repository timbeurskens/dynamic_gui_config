package dynamic_gui_config

import "github.com/andlabs/ui"

type horizontalValueControlArray []ValueControl

func (h horizontalValueControlArray) Create() ui.Control {
	hbox := ui.NewHorizontalBox()

	for _, control := range h {
		hbox.Append(control.Create(), true)
	}

	return hbox
}

type verticalValueControlArray []ValueControl

func (v verticalValueControlArray) Create() ui.Control {
	vbox := ui.NewVerticalBox()

	for _, control := range v {
		vbox.Append(control.Create(), true)
	}

	return vbox
}
