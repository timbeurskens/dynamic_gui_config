package dynamic_gui_config

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

var (
	valueControlType = reflect.TypeOf((*ValueControl)(nil)).Elem()
	onChangedType    = reflect.TypeOf((*UpdateNotifier)(nil)).Elem()
)

// new

func MakeValueControl(object interface{}) (ValueControl, error) {
	return MakeValueControlFromValue(reflect.ValueOf(object), structTagDefaults)
}

func MakeValueControlFromValue(value reflect.Value, properties StructTagProperties) (ValueControl, error) {
	if !value.IsValid() {
		return nil, errors.New("invalid value")
	}

	if !value.CanInterface() {
		return nil, errors.New("cannot interface value")
	}

	if value.Type().Implements(valueControlType) {
		return value.Interface().(ValueControl), nil
	}

	switch value.Kind() {
	case reflect.Ptr:
		return MakeValueControlFromValue(value.Elem(), properties)
	case reflect.Struct:
		return structBreakdown(value, properties)
	case reflect.Array, reflect.Slice:
		return arrayBreakdown(value, properties)
	default:
		return fieldValueBreakdown(value, properties)
	}
}

// old

func fieldValueBreakdown(value reflect.Value, properties StructTagProperties) (ValueControl, error) {
	if !value.CanAddr() {
		return nil, errors.New("cannot take address of value")
	}

	if bIn, ok := builtin[value.Kind()]; ok {
		log.Printf("adding builtin object %s kind %s", value.Type(), value.Kind())
		onchanged := func() {}

		if value.Addr().Type().Implements(onChangedType) {
			iface, ok := value.Addr().Interface().(UpdateNotifier)
			if ok {
				log.Printf("add onchanged handler for type %s", value.Type())
				onchanged = iface.OnValueChanged
			}
		}

		kindType, ok := typedKind[value.Kind()]

		if ok && value.Addr().Type().ConvertibleTo(kindType) {
			valueConverted := value.Addr().Convert(typedKind[value.Kind()])

			if controlFactory := bIn(valueConverted.Interface(), properties, onchanged); controlFactory != nil {
				return controlFactory, nil
			} else {
				return nil, errors.New(fmt.Sprintf("cannot create builtin object %s", value.Kind()))
			}
		} else {
			return nil, errors.New(fmt.Sprintf("cannot convert %s to %s", value.Type(), kindType))
		}
	}
	return nil, errors.New(fmt.Sprintf("no match for %s", value.Type()))
}

func arrayBreakdown(array reflect.Value, properties StructTagProperties) (ValueControl, error) {
	result := make([]ValueControl, 0, array.Len())

	for i := 0; i < array.Len(); i++ {
		if valueBreakdown, err := MakeValueControlFromValue(array.Index(i), properties); err != nil {
			log.Printf("ignoring %s[%d]: %s", array.Type(), i, err)
		} else {
			result = append(result, valueBreakdown)
		}
	}

	if properties.Horizontal {
		return horizontalValueControlArray(result), nil
	} else {
		return verticalValueControlArray(result), nil
	}
}

func fieldBreakdown(field reflect.Value, structField reflect.StructField) (ValueControl, error) {
	properties, err := ParseStructTag(structField.Tag.Get(tagkey))
	if err != nil {
		log.Printf("error parsing struct tag: %s", err)
	}

	if properties.Name == "" {
		properties.Name = Label(structField.Name)
	}

	if factory, err := MakeValueControlFromValue(field, properties); err != nil {
		return nil, err
	} else {
		return structGuiField{
			Properties: properties,
			Factory:    factory,
		}, nil
	}
}

func structBreakdown(value reflect.Value, properties StructTagProperties) (ValueControl, error) {
	result := make([]ValueControl, 0, value.NumField())

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)

		if guifield, err := fieldBreakdown(field, structField); err == nil {
			result = append(result, guifield)
		} else {
			log.Printf("Ignoring %s, error: %s, ", structField.Name, err)
		}
	}

	if properties.Horizontal {
		return horizontalValueControlArray(result), nil
	} else {
		return verticalValueControlArray(result), nil
	}
}
