package dynamic_gui_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/andlabs/ui"
)

var (
	valueControlType = reflect.TypeOf((*ValueControl)(nil)).Elem()
	onChangedType    = reflect.TypeOf((*UpdateNotifier)(nil)).Elem()
)

const (
	tagkey = "uiconf"
)

type ValueControlFunc func() ui.Control

func (v ValueControlFunc) Create() ui.Control {
	return v()
}

// UpdateNotifier is an interface type a user can implement to get notified when a value changes
type UpdateNotifier interface {
	// OnValueChanged will be called when the user changes the value using the graphical interface
	OnValueChanged()
}

type ValueControl interface {
	Create() ui.Control
}

type horizontalValueControlArray []ValueControl

func (h horizontalValueControlArray) Create() ui.Control {
	hbox := ui.NewHorizontalBox()

	for _, control := range h {
		hbox.Append(control.Create(), true)
	}

	return hbox
}

type StructTagProperties struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Min        int    `json:"min,omitempty"`
	Max        int    `json:"max,omitempty"`
	Resolution int    `json:"resolution,omitempty"`
}

type structGuiField struct {
	Properties StructTagProperties
	Factory    ValueControl
}

func parseStructTag(tag string) (properties StructTagProperties, err error) {
	properties = StructTagProperties{
		Name:       "",
		Type:       "",
		Min:        0,
		Max:        100,
		Resolution: 10,
	}

	if tag == "" {
		err = errors.New("empty tag, using defaults")
		return
	}

	err = json.Unmarshal([]byte(tag), &properties)
	return
}

func fieldValueBreakdown(value reflect.Value, properties StructTagProperties) (ValueControl, error) {
	if !value.CanAddr() {
		return nil, errors.New("cannot take address of value")
	}
	fieldAddr := value.Addr()

	if !fieldAddr.CanInterface() {
		return nil, errors.New("cannot take interface of value pointer")
	}

	if fieldAddr.IsNil() {
		return nil, errors.New("address of value is nil")
	}

	fieldAddrIface := fieldAddr.Interface()

	if value.Type().Implements(valueControlType) {
		log.Printf("adding ValueControl object %s implementing %s", value.Type(), valueControlType)
		return fieldAddrIface.(ValueControl), nil
	} else if value.Kind() == reflect.Array {
		if controlFactory, err := arrayBreakdown(value, properties); err != nil {
			return nil, err
		} else {
			return controlFactory, nil
		}
	} else if _, ok := builtin[value.Kind()]; ok {
		log.Printf("adding builtin object %s kind %s", value.Type(), value.Kind())
		onchanged := func() {}

		if value.Type().Implements(onChangedType) {
			iface, ok := fieldAddr.Interface().(UpdateNotifier)
			if ok {
				log.Printf("add onchanged handler for type %s", value.Type())
				onchanged = iface.OnValueChanged
			}
		}

		kindType, ok := typedKind[value.Kind()]

		if ok && fieldAddr.Type().ConvertibleTo(kindType) {
			valueConverted := fieldAddr.Convert(typedKind[value.Kind()])

			if controlFactory := builtin[value.Kind()](valueConverted.Interface(), properties, onchanged); controlFactory != nil {
				return controlFactory, nil
			}
		} else {
			return nil, errors.New(fmt.Sprintf("cannot convert %s to %s", fieldAddr.Type(), kindType))
		}
	}
	return nil, errors.New(fmt.Sprintf("no match for %s", value.Type()))
}

func arrayBreakdown(array reflect.Value, properties StructTagProperties) (ValueControl, error) {
	result := make(horizontalValueControlArray, 0)

	for i := 0; i < array.Len(); i++ {
		if valueBreakdown, err := fieldValueBreakdown(array.Index(i), properties); err != nil {
			log.Printf("ignoring %s[%d]: %s", array.Type(), i, err)
		} else {
			result = append(result, valueBreakdown)
		}
	}

	return result, nil
}

func fieldBreakdown(field reflect.Value, structField reflect.StructField) (structGuiField, error) {
	properties, err := parseStructTag(structField.Tag.Get(tagkey))
	if err != nil {
		log.Printf("error parsing struct tag: %s", err)
	}

	if properties.Name == "" {
		properties.Name = structField.Name
	}

	if factory, err := fieldValueBreakdown(field, properties); err != nil {
		return structGuiField{}, err
	} else {
		return structGuiField{
			Properties: properties,
			Factory:    factory,
		}, nil
	}
}

func structBreakdown(structPtr interface{}) ([]structGuiField, error) {
	reflectValue := reflect.ValueOf(structPtr)

	if reflectValue.Kind() != reflect.Ptr {
		return nil, errors.New(fmt.Sprintf("structPtr should be a pointer to a struct type, got %d", reflectValue.Kind()))
	}

	// should be safe now, already checked for pointer type
	value := reflectValue.Elem()

	if value.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("structPtr should be a pointer to a struct type, got pointer to %s", value.Kind()))
	}

	result := make([]structGuiField, 0)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)

		if guifield, err := fieldBreakdown(field, structField); err == nil {
			result = append(result, guifield)
		} else {
			log.Printf("Ignoring %s, error: %s, ", structField.Name, err)
		}
	}

	return result, nil
}
