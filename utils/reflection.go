package utils

import (
	"reflect"
)

func Inspect(data any) map[string]any {
	ret := make(map[string]any, 0)
	val := reflect.ValueOf(data).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		f := valueField.Interface()
		val := reflect.ValueOf(f)

		if v := convertFieldToType(val); v != nil {
			ret[typeField.Name] = v
		}
	}

	return ret
}

func convertFieldToType(f reflect.Value) any {
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int()
	case reflect.Float32, reflect.Float64:
		return f.Float()
	case reflect.String:
		return f.String()
	case reflect.Bool:
		return f.Bool()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Uint()
	default:
		return nil
	}
}
