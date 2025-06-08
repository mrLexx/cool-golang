package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

const ValidateTag = "validate"

var rulesStore = ListRules{
	"len": {
		Type:         []reflect.Kind{reflect.String},
		ValidateData: ValidateLen,
	},
	"regexp": {
		Type:         []reflect.Kind{reflect.String},
		ValidateData: ValidateRegexp,
	},
	"in": {
		Type:         []reflect.Kind{reflect.String, reflect.Int},
		ValidateData: ValidateIn,
	},
	"out": {
		Type:         []reflect.Kind{reflect.String, reflect.Int},
		ValidateData: ValidateOut,
	},
	"min": {
		Type:         []reflect.Kind{reflect.Int},
		ValidateData: ValidateMin,
	},
	"max": {
		Type:         []reflect.Kind{reflect.Int},
		ValidateData: ValidateMax,
	},
}

var validErrs = make(ValidationErrors, 0)

func Validate(v any) error {
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
		if err := validateObj(val); err != nil {
			return err
		}

	default:
		if err := validateObj(val); err != nil {
			return err
		}
	}
	fmt.Println("len valid Errors: ", len(validErrs))
	for _, v := range validErrs {
		fmt.Println(v.Error())

	}
	return nil
}

func validateObj(val reflect.Value) error {

	if val.Kind() == reflect.Slice {
		for i := range val.Len() {
			elem := val.Index(i)
			if elem.Kind() == reflect.Pointer {
				if elem.IsNil() {
					continue
				}
				elem = elem.Elem()
			}
			err := validateItem(elem)
			var validErr *ValidationError
			switch {
			case errors.As(err, &validErr):
				validErrs = append(validErrs, *validErr)
			default:
				return err
			}

		}
	} else {
		err := validateItem(val)
		var validErr *ValidationError
		switch {
		case errors.As(err, &validErr):
			validErrs = append(validErrs, *validErr)
		default:
			return err
		}
	}

	return nil
}

func validateItem(v reflect.Value) error {

	vt := v.Type()

	for i := range vt.NumField() {
		f := vt.Field(i)

		if tag, ok := f.Tag.Lookup(ValidateTag); ok {
			if tag == "nested" {
				vn := v.Field(i)
				if err := separateValidateError(validateObj(vn)); err != nil {
					return err
				}
			} else {
				for _, t := range splitTag(tag) {
					rs, err := extractRule(t)
					if err != nil {
						return err
					}
					itm, ok := rulesStore[rs.Name]
					if !ok {
						return NewExecuteError(ErrExecuteUndefinedRule, "has an undefined rule `%v`", rs.Name)
					}

					if err := separateValidateError(itm.ValidateData(rs.Payload, v.Field(i).Interface())); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
