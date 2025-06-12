package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const ValidateTag = "validate"

type ruleSet struct {
	Name    string
	Payload string
}

var rulesStore = ListRules{
	"len": {
		Validate: validateLen,
	},
	"regexp": {
		Validate: validateRegexp,
	},
	"require": {
		Validate: validateRequire,
	},
	"in": {
		Validate: validateIn,
	},
	"out": {
		Validate: validateOut,
	},
	"min": {
		Validate: validateMin,
	},
	"max": {
		Validate: validateMax,
	},
}

var regexpCache = NewCache[*regexp.Regexp]()

var inCache = NewCache[bool]()

var outCache = NewCache[bool]()

func Validate(v any) error {
	validationErrs := make(ValidationErrors, 0)

	val := reflect.ValueOf(v)
	typ := val.Type()
	typKind := typ.Kind()

	switch {
	case typKind != reflect.Struct && typKind != reflect.Slice:
		return NewExecuteError(ErrExecuteWrongInput, "expected struct or slice of structs, got %T", v)

	case typKind == reflect.Slice:
		elemKind := typ.Elem().Kind()
		if elemKind != reflect.Struct && !(elemKind == reflect.Pointer && typ.Elem().Elem().Kind() == reflect.Struct) {
			return NewExecuteError(ErrExecuteWrongInput, "expected slice of structs or *structs, got %T", v)
		}
		if err := validateStruct(val, &validationErrs); err != nil {
			return err
		}

	default:
		if err := validateStruct(val, &validationErrs); err != nil {
			return err
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func validateStruct(val reflect.Value, validationErrs *ValidationErrors) error {
	switch {
	case val.Kind() == reflect.Slice:
		for i := range val.Len() {
			elem := val.Index(i)
			if elem.Kind() == reflect.Pointer {
				if elem.IsNil() {
					continue
				}
				elem = elem.Elem()
			}
			err := validateField(elem, validationErrs)
			if err != nil {
				return err
			}
		}
	default:
		err := validateField(val, validationErrs)
		if err != nil {
			return fmt.Errorf("%v.%w", val.Type().Name(), err)
		}
	}
	return nil
}

func validateField(v reflect.Value, validationErrs *ValidationErrors) error {
	vt := v.Type()
	for i := range vt.NumField() {
		f := vt.Field(i)

		tag, ok := f.Tag.Lookup(ValidateTag)
		if !ok {
			continue
		}

		vn := v.Field(i)

		switch {
		case tag == "nested":
			if err := validateStruct(vn, validationErrs); err != nil {
				return fmt.Errorf("%v: %w", f.Name, err)
			}
		default:
			if err := validateData(f.Name, tag, vn.Interface(), validationErrs); err != nil {
				// return err
				return fmt.Errorf("%v: %w", f.Name, err)
			}
		}
	}
	return nil
}

func validateData(fName, tag string, v any, validationErrs *ValidationErrors) error {
	for _, r := range splitRules(tag) {
		tp := reflect.TypeOf(v)
		kn := tp.Kind()

		rS, err := extractRule(r)
		if err != nil {
			return err
		}
		itm, ok := rulesStore[rS.Name]
		if !ok {
			return NewExecuteError(ErrExecuteUndefinedRule, "has an undefined rule `%v`", rS.Name)
		}

		vRef := reflect.ValueOf(v)
		switch {
		case kn == reflect.Slice:
			kn := tp.Elem().Kind()
			if vRef.Len() == 0 {
				vRef = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(vRef)), 1, 1)
			}
			for i := range vRef.Len() {
				if err := separateValidationError(
					itm.Validate(rS.Payload, valueSet{Val: vRef.Index(i), Type: kn}),
					fName,
					validationErrs,
				); err != nil {
					return err
				}
			}
		default:
			if err := separateValidationError(
				itm.Validate(rS.Payload, valueSet{Val: vRef, Type: kn}),
				fName,
				validationErrs,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractRule(rawRule string) (ruleSet, error) {
	if rawRule == "require" {
		rawRule = "require:true"
	}
	tmp := strings.Split(rawRule, ":")
	if len(tmp) != 2 || tmp[1] == "" {
		return ruleSet{}, NewExecuteError(ErrExecuteIncompleteRule, "has an incomplete rule `%v`", rawRule)
	}
	r, p := tmp[0], tmp[1]
	return ruleSet{Name: r, Payload: p}, nil
}

func splitRules(tag string) []string {
	return strings.Split(tag, "|")
}
