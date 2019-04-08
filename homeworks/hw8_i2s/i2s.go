package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	outValue := reflect.ValueOf(out)

	if outValue.Kind() != reflect.Ptr {
		return fmt.Errorf("out is not a pointer")
	} else {
		outValue = outValue.Elem() // разыменовываем указатель
	}

	switch outValue.Kind() {
	case reflect.Struct:
		d, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed convert to map[string]interface{}")
		}

		for i := 0; i < outValue.NumField(); i++ {
			fieldName := outValue.Type().Field(i).Name

			v, ok := d[fieldName]
			if !ok {
				return fmt.Errorf("field not found: %s", fieldName)
			}

			if err := i2s(v, outValue.Field(i).Addr().Interface()); err != nil {
				return fmt.Errorf("failed to process struct field %s: %s", fieldName, err)
			}
		}

	case reflect.Slice:
		d, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("failed convert to []interface{}")
		}

		for i, v := range d {
			o := reflect.New(outValue.Type().Elem())
			if err := i2s(v, o.Interface()); err != nil {
				return fmt.Errorf("failed to process slice element %d: %s", i, err)
			}

			outValue.Set(reflect.Append(outValue, o.Elem()))
		}

	case reflect.Int:
		d, ok := data.(float64) // json распаковывает int во float
		if !ok {
			return fmt.Errorf("failed convert to float64")
		}

		outValue.SetInt(int64(d))

	case reflect.String:
		d, ok := data.(string)
		if !ok {
			return fmt.Errorf("failed convert to string")
		}

		outValue.SetString(d)

	case reflect.Bool:
		d, ok := data.(bool)
		if !ok {
			return fmt.Errorf("failed convert to bool")
		}

		outValue.SetBool(d)

	default:
		return fmt.Errorf("unsupportd type: %s", outValue.Kind())
	}

	return nil
}
