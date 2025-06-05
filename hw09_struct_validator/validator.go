package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

var regexpList = make(map[string]*regexp.Regexp, 0)

func Validate(v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return makeExecuteErrorf(ErrExecuteWrongInput, "expected a struct, but received %T", v)
	}

	if err := validateStruct(reflect.ValueOf(v).Type()); err != nil {
		return err
	}

	return nil
}

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

	for _, v := range strings.Split(tag, "|") {
		r, p, err := extractRule(v)
		if err != nil {
			return fmt.Errorf("field `%v`: %w", sf.Name, err)
		}

		if err := validateMappingFieldRule(fieldType, rulesStore[r].Type); err != nil {
			return fmt.Errorf("field `%v` (%v): rule `%v`: %w", sf.Name, fieldType, r, err)
		}

		if err := validateRule(r, p, fieldType); err != nil {
			return fmt.Errorf("field `%v`: %w", sf.Name, err)
		}
	}

	return nil
}

func validateMappingFieldRule(tp reflect.Kind, expectedTp []reflect.Kind) error {
	if !slices.Contains(expectedTp, tp) {
		return makeExecuteErrorf(ErrExecuteWrongRuleType, "this rule only for (%v)", expectedTp)
	}
	return nil
}

func validateRule(r, p string, fieldType reflect.Kind) error {
	if _, ok := rulesStore[r]; !ok {
		return makeExecuteErrorf(ErrExecuteUndefinedRule, "has an undefined rule `%v`", r)
	}
	return rulesStore[r].Validate(r, p, fieldType)
}
