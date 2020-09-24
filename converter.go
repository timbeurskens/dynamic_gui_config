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
	fieldAddr := value.Addr()

	if !fieldAddr.CanInterface() {
		return nil, errors.New("cannot take interface of value pointer")
	}

	if fieldAddr.IsNil() {
		return nil, errors.New("address of value is nil")
	}

	if value.Type().Implements(valueControlType) {
		log.Printf("adding ValueControl object %s implementing %s", value.Type(), valueControlType)
		return fieldAddr.Interface().(ValueControl), nil
	} else if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		return arrayBreakdown(value, properties)
	} else if value.Kind() == reflect.Struct {
		return structBreakdown(fieldAddr, properties)
	} else if bIn, ok := builtin[value.Kind()]; ok {
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

			if controlFactory := bIn(valueConverted.Interface(), properties, onchanged); controlFactory != nil {
				return controlFactory, nil
			} else {
				return nil, errors.New(fmt.Sprintf("cannot create builtin object %s", value.Kind()))
			}
		} else {
			return nil, errors.New(fmt.Sprintf("cannot convert %s to %s", fieldAddr.Type(), kindType))
		}
	}
	return nil, errors.New(fmt.Sprintf("no match for %s", value.Type()))
}

func arrayBreakdown(array reflect.Value, properties StructTagProperties) (ValueControl, error) {
	result := make([]ValueControl, 0)

	for i := 0; i < array.Len(); i++ {
		if valueBreakdown, err := fieldValueBreakdown(array.Index(i), properties); err != nil {
			log.Printf("ignoring %s[%d]: %s", array.Type(), i, err)
		} else {
			result = append(result, valueBreakdown)
		}
	}

	if properties.Vertical {
		return verticalValueControlArray(result), nil
	} else {
		return horizontalValueControlArray(result), nil
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

	if factory, err := fieldValueBreakdown(field, properties); err != nil {
		return nil, err
	} else {
		return structGuiField{
			Properties: properties,
			Factory:    factory,
		}, nil
	}
}

func structBreakdown(reflectValue reflect.Value, properties StructTagProperties) (ValueControl, error) {
	if reflectValue.Kind() != reflect.Ptr {
		return nil, errors.New(fmt.Sprintf("structPtr should be a pointer to a struct type, got %d", reflectValue.Kind()))
	}

	// should be safe now, already checked for pointer type
	value := reflectValue.Elem()

	if value.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("structPtr should be a pointer to a struct type, got pointer to %s", value.Kind()))
	}

	result := make(controlGroup, 0, value.NumField())
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
