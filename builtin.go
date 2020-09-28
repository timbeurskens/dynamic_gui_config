package dynamic_gui_config

import (
	"errors"
	"reflect"

	"github.com/andlabs/ui"
)

// map kinds to pointer types
var typedKind = map[reflect.Kind]reflect.Type{
	reflect.Float64: reflect.TypeOf((*float64)(nil)),
	reflect.Int:     reflect.TypeOf((*int)(nil)),
	reflect.Bool:    reflect.TypeOf((*bool)(nil)),
	reflect.String:  reflect.TypeOf((*string)(nil)),
	reflect.Float32: reflect.TypeOf((*float32)(nil)),
	reflect.Uint:    reflect.TypeOf((*uint)(nil)),
	reflect.Func:    reflect.TypeOf((*func())(nil)),
}

var builtin = map[reflect.Kind]func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error){
	reflect.Float64: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		value, ok := i.(*float64)
		if !ok {
			return nil, errors.New("could not convert i to float64")
		}

		return ValueControlFunc(func() ui.Control {
			slider := ui.NewSlider(properties.Min*properties.Resolution, properties.Max*properties.Resolution)
			slider.SetValue(int(*value * float64(properties.Resolution)))
			slider.OnChanged(func(slider *ui.Slider) {
				*value = float64(slider.Value()) / float64(properties.Resolution)
				go onchanged()
			})
			return slider
		}), nil
	},
	reflect.Float32: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		value, ok := i.(*float32)
		if !ok {
			return nil, errors.New("could not convert i to float32")
		}

		return ValueControlFunc(func() ui.Control {
			slider := ui.NewSlider(properties.Min*properties.Resolution, properties.Max*properties.Resolution)
			slider.SetValue(int(*value * float32(properties.Resolution)))
			slider.OnChanged(func(slider *ui.Slider) {
				*value = float32(slider.Value()) / float32(properties.Resolution)
				go onchanged()
			})
			return slider
		}), nil
	},
	reflect.Uint: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		value, ok := i.(*uint)
		if !ok {
			return nil, errors.New("could not convert i to uint")
		}

		return ValueControlFunc(func() ui.Control {
			slider := ui.NewSlider(0, properties.Max)
			slider.SetValue(int(*value))
			slider.OnChanged(func(slider *ui.Slider) {
				*value = uint(slider.Value())
				go onchanged()
			})
			return slider
		}), nil
	},
	reflect.Int: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		value, ok := i.(*int)
		if !ok {
			return nil, errors.New("could not convert i to int")
		}

		return ValueControlFunc(func() ui.Control {
			slider := ui.NewSlider(properties.Min, properties.Max)
			slider.SetValue(*value)
			slider.OnChanged(func(slider *ui.Slider) {
				*value = slider.Value()
				go onchanged()
			})
			return slider
		}), nil
	},
	reflect.Func: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		callbackPtr, ok := i.(*func())
		if !ok {
			return nil, errors.New("could not convert i to func")
		}

		callback := *callbackPtr

		// extra check for nil function pointers
		if callback == nil {
			return nil, errors.New("callback refers to nil")
		}

		return ValueControlFunc(func() ui.Control {
			button := ui.NewButton(string(properties.Name))
			button.OnClicked(func(button *ui.Button) {
				go callback()
			})
			return button
		}), nil
	},
	reflect.String: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		value, ok := i.(*string)
		if !ok {
			return nil, errors.New("could not convert i to string")
		}

		return ValueControlFunc(func() ui.Control {
			textbox := ui.NewEntry()
			textbox.SetText(*value)
			textbox.OnChanged(func(entry *ui.Entry) {
				*value = entry.Text()
				go onchanged()
			})
			return textbox
		}), nil
	},
	reflect.Bool: func(i interface{}, properties StructTagProperties, onchanged func()) (ValueControl, error) {
		value, ok := i.(*bool)
		if !ok {
			return nil, errors.New("could not convert i to bool")
		}

		return ValueControlFunc(func() ui.Control {
			check := ui.NewCheckbox("")
			check.SetChecked(*value)
			check.OnToggled(func(checkbox *ui.Checkbox) {
				*value = checkbox.Checked()
				go onchanged()
			})
			return check
		}), nil
	},
}
