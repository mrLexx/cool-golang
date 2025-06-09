package hw09structvalidator

import "reflect"

func getFieldType(f reflect.StructField) reflect.Kind {
	if f.Type.Kind() == reflect.Slice {
		return f.Type.Elem().Kind()
	}
	return f.Type.Kind()
}
