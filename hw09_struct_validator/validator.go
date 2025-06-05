package hw09structvalidator

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func Validate(v any) error {

	// validateStore := make(map[string][]map[string]string, 0)

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return makeExecuteErrorf(ErrExecuteWrongInput, "expected a struct, but received %T", v)
	}

	if err := validateStruct(reflect.ValueOf(v)); err != nil {
		return err
	}
	// vt := reflect.TypeOf(v)

	// extract rules - start
	// for i := range vt.NumField() {
	// 	f := vt.Field(i)
	// 	if tag, ok := f.Tag.Lookup(ValidateTag); ok {
	// 		tt, err := exctractTag(f, tag)
	// 		if err != nil {
	// 			return fmt.Errorf("type `%v`: %w", vt.Name(), err)
	// 		}
	// 		validateStore[f.Name] = tt
	// 	}
	// }
	// fmt.Println(validateStore)
	// extract rules

	return nil
}

func validateStruct(v reflect.Value) error {
	vt := v.Type()

	for i := range vt.NumField() {
		f := vt.Field(i)
		if tag, ok := f.Tag.Lookup(ValidateTag); ok {
			_, err := exctractTag(f, tag)
			if err != nil {
				return fmt.Errorf("type `%v`: %w", vt.Name(), err)
			}
			// validateStore[f.Name] = tt
		}
	}

	return nil
}

func exctractTag(f reflect.StructField, tag string) ([]map[string]string, error) {
	rules := make([]map[string]string, 0, strings.Count(tag, "|")+1)

	if tag == "nested" {
		return nil, nil
	}

	fieldType := getFieldType(f)

	for _, v := range strings.Split(tag, "|") {
		r, p, err := parseTag(v)
		if err != nil {
			return nil, fmt.Errorf("field `%v`: %w", f.Name, err)
		}

		if err := checkMappingFieldRule(fieldType, rulesStore[r].Type); err != nil {
			return nil, fmt.Errorf("field `%v` (%v): rule `%v`: %w", f.Name, fieldType, r, err)
		}

		if err := prepareRule(r, p, fieldType); err != nil {
			return nil, fmt.Errorf("field `%v`: %w", f.Name, err)
		}

		rules = append(rules, map[string]string{r: p})
	}

	return rules, nil
}

func parseTag(tag string) (r, p string, err error) {
	tmp := strings.Split(tag, ":")

	if len(tmp) != 2 || tmp[1] == "" {
		fmt.Println(tmp[0])
		return "", "", makeExecuteErrorf(ErrExecuteIncompleteRule, "has an incomplete rule `%v`", tag)
	}

	r, p = tmp[0], tmp[1]

	if _, ok := rulesStore[r]; !ok {
		return "", "", makeExecuteErrorf(ErrExecuteUndefinedRule, "has an undefined rule `%v`", r)
	}

	return
}

func checkMappingFieldRule(tp reflect.Kind, expectedTp []reflect.Kind) error {
	if !slices.Contains(expectedTp, tp) {
		return makeExecuteErrorf(ErrExecuteWrongRuleType, "this rule only for (%v)", expectedTp)
	}
	return nil
}

func prepareRule(r, p string, fieldType reflect.Kind) error {
	if _, ok := rulesStore[r]; !ok {
		return makeExecuteErrorf(ErrExecuteUndefinedRule, "has an undefined rule `%v`", r)
	}
	return rulesStore[r].Prepare(r, p, fieldType)
}
