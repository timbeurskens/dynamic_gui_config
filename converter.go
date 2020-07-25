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

func fieldBreakdown(field reflect.Value, structField reflect.StructField) (structGuiField, error) {
	if !field.CanAddr() {
		return structGuiField{}, errors.New("cannot take address of field")
	}
	fieldAddr := field.Addr()

	if !fieldAddr.CanInterface() {
		return structGuiField{}, errors.New("cannot take interface of field pointer")
	}

	fieldAddrIface := fieldAddr.Interface()

	properties, err := parseStructTag(structField.Tag.Get(tagkey))
	if err != nil {
		log.Printf("error parsing struct tag: %s", err)
	}

	if properties.Name == "" {
		properties.Name = structField.Name
	}

	log.Println(field.Kind().String())

	log.Println(field.Type().String())

	if field.Type().Implements(valueControlType) {
		log.Printf("adding ValueControl object %s implementing %s", field.Type(), valueControlType)
		return structGuiField{
			Properties: properties,
			Factory:    fieldAddrIface.(ValueControl),
		}, nil
	} else if _, ok := builtin[field.Kind()]; ok {
		log.Printf("adding builtin object %s kind %s", field.Type(), field.Kind())
		onchanged := func() {}

		if field.Type().Implements(onChangedType) {
			iface, ok := fieldAddr.Interface().(UpdateNotifier)
			if ok {
				log.Printf("add onchanged handler for type %s", field.Type())
				onchanged = iface.OnValueChanged
			}
		}

		return structGuiField{
			Properties: properties,
			Factory:    builtin[field.Kind()](fieldAddr.Convert(typedKind[field.Kind()]).Interface(), properties, onchanged),
		}, nil
	}

	return structGuiField{
		Properties: StructTagProperties{},
		Factory:    nil,
	}, errors.New(fmt.Sprintf("no match for %s", field.Type()))
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
