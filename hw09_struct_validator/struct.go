package hw09structvalidator

import (
	"fmt"
	"reflect"
	"slices"
)

func validateStruct(t reflect.Type) error {
	for i := range t.NumField() {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup(ValidateTag); ok {
			if err := validateTag(f, tag); err != nil {
				return fmt.Errorf("type `%v`: %w", t.Name(), err)
			}
		}
	}
	return nil
}

func validateTag(sf reflect.StructField, tag string) error {
	if tag == "nested" {
		return validateStruct(sf.Type)
	}

	fieldType := getFieldType(sf)

	for _, v := range splitTag(tag) {
		// r, p, err := extractRule(v)
		rr, err := extractRule(v)
		if err != nil {
			return fmt.Errorf("field `%v`: %w", sf.Name, err)
		}

		if err := validateMappingFieldRule(fieldType, rulesStore[rr.Name].Type); err != nil {
			return fmt.Errorf("field `%v` (%v): rule `%v`: %w", sf.Name, fieldType, rr.Name, err)
		}

		if err := validateRule(rr.Name, rr.Payload, fieldType); err != nil {
			return fmt.Errorf("field `%v`: %w", sf.Name, err)
		}
	}

	return nil
}

func validateMappingFieldRule(tp reflect.Kind, expectedTp []reflect.Kind) error {
	if !slices.Contains(expectedTp, tp) {
		return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (%v)", expectedTp)
	}
	return nil
}

func validateRule(r, p string, fieldType reflect.Kind) error {
	if _, ok := rulesStore[r]; !ok {
		return NewExecuteError(ErrExecuteUndefinedRule, "has an undefined rule `%v`", r)
	}
	return rulesStore[r].ValidateStruct(p, fieldType)
}
