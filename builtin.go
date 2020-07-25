package dynamic_gui_config

import (
	"reflect"

	"github.com/andlabs/ui"
)

// map kinds to pointer types
var typedKind = map[reflect.Kind]reflect.Type{
	reflect.Float64: reflect.TypeOf((*float64)(nil)),
	reflect.Int:     reflect.TypeOf((*int)(nil)),
	reflect.Bool:    reflect.TypeOf((*bool)(nil)),
	reflect.String:  reflect.TypeOf((*string)(nil)),
}

var builtin = map[reflect.Kind]func(i interface{}, properties StructTagProperties, onchanged func()) ValueControl{
	reflect.Float64: func(i interface{}, properties StructTagProperties, onchanged func()) ValueControl {
		value := i.(*float64)

		return ValueControlFunc(func() ui.Control {
			slider := ui.NewSlider(properties.Min*properties.Resolution, properties.Max*properties.Resolution)
			slider.SetValue(int(*value * float64(properties.Resolution)))
			slider.OnChanged(func(slider *ui.Slider) {
				*value = float64(slider.Value()) / float64(properties.Resolution)
				go onchanged()
			})
			return slider
		})
	},
	reflect.Int: func(i interface{}, properties StructTagProperties, onchanged func()) ValueControl {
		value := i.(*int)

		return ValueControlFunc(func() ui.Control {
			slider := ui.NewSlider(properties.Min, properties.Max)
			slider.SetValue(*value)
			slider.OnChanged(func(slider *ui.Slider) {
				*value = slider.Value()
				go onchanged()
			})
			return slider
		})
	},
	reflect.String: func(i interface{}, properties StructTagProperties, onchanged func()) ValueControl {
		value := i.(*string)

		return ValueControlFunc(func() ui.Control {
			textbox := ui.NewEntry()
			textbox.SetText(*value)
			textbox.OnChanged(func(entry *ui.Entry) {
				*value = entry.Text()
				go onchanged()
			})
			return textbox
		})
	},
	reflect.Bool: func(i interface{}, properties StructTagProperties, onchanged func()) ValueControl {
		value := i.(*bool)

		return ValueControlFunc(func() ui.Control {
			check := ui.NewCheckbox("")
			check.SetChecked(*value)
			check.OnToggled(func(checkbox *ui.Checkbox) {
				*value = checkbox.Checked()
				go onchanged()
			})
			return check
		})
	},
}
