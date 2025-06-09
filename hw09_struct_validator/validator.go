package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
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

var validationErrs = make(ValidationErrors, 0)

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
		if err := validateStruct(val); err != nil {
			return err
		}

	default:
		if err := validateStruct(val); err != nil {
			return err
		}
	}
	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func validateStruct(val reflect.Value) error {
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
			err := validateField(elem)
			if err != nil {
				return err
			}
		}
	default:
		err := validateField(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateField(v reflect.Value) error {
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
			if err := validateStruct(vn); err != nil {
				return err
			}
		default:
			if err := validateTag(f.Name, tag, vn.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateTag(fName, tag string, v any) error {
	for _, r := range splitRules(tag) {
		rs, err := extractRule(r)
		if err != nil {
			return err
		}
		itm, ok := rulesStore[rs.Name]
		if !ok {
			return NewExecuteError(ErrExecuteUndefinedRule, "has an undefined rule `%v`", rs.Name)
		}

		err = itm.ValidateData(rs.Payload, v)

		var execErr *ExecuteError
		if errors.As(err, &execErr) {
			return err
		}
		if err != nil {
			validationErrs = append(validationErrs, ValidationError{
				Field: fName,
				Err:   err,
			})
		}
	}
	return nil
}

func extractRule(tag string) (RuleSet, error) {
	tmp := strings.Split(tag, ":")
	if len(tmp) != 2 || tmp[1] == "" {
		return RuleSet{}, NewExecuteError(ErrExecuteIncompleteRule, "has an incomplete rule `%v`", tag)
	}
	r, p := tmp[0], tmp[1]
	return RuleSet{Name: r, Payload: p}, nil
}

func splitRules(tag string) []string {
	return strings.Split(tag, "|")
}
