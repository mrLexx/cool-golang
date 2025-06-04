package hw09structvalidator

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func Validate(v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return makeExecuteErrorf(ErrExecuteWrongInput, "expected a struct, but received %T", v)
	}

	vt := reflect.TypeOf(v)

	// extract rules - start
	for i := range vt.NumField() {
		f := vt.Field(i)
		if tag, ok := f.Tag.Lookup(ValidateTag); ok {
			_, err := exctractRules(f, tag)
			if err != nil {
				return fmt.Errorf("type `%v`: %w", vt.Name(), err)
			}
		}
	}
	// extract rules

	return nil
}

func exctractRules(f reflect.StructField, tag string) ([]map[string]string, error) {
	rules := make([]map[string]string, 0, strings.Count(tag, "|")+1)

	if tag == "nested" {
		return rules, nil
	}

	for _, v := range strings.Split(tag, "|") {
		r, p, err := parseTag(f, v)
		if err != nil {
			return nil, err
		}

		if err := prepareRule(f, r, p); err != nil {
			return nil, err
		}

		rules = append(rules, map[string]string{r: p})
	}

	return rules, nil
}

func parseTag(f reflect.StructField, tag string) (r, p string, err error) {
	tmp := strings.Split(tag, ":")

	if len(tmp) != 2 || tmp[1] == "" {
		fmt.Println(tmp[0])
		return "", "", makeExecuteErrorf(ErrExecuteIncompleteRule, "field `%v` has an incomplete `%v` rule", f.Name, tag)
	}

	r, p = tmp[0], tmp[1]

	if _, ok := rulesList[r]; !ok {
		return "", "", makeExecuteErrorf(ErrExecuteUndefinedRule, "field `%v` has an undefined rule `%v`", f.Name, r)
	}

	fType := getReflectType(f)

	if !slices.Contains(rulesList[r].Type, fType) {
		return "", "", makeExecuteErrorf(ErrExecuteWrongRuleType,
			"field `%v` (%v) has wrong rule type `%v` (%v)", f.Name, fType, r, rulesList[tmp[0]].Type)
	}

	return
}

func prepareRule(f reflect.StructField, r, p string) error {
	if _, ok := rulesList[r]; !ok {
		return makeExecuteErrorf(ErrExecuteUndefinedRule, "field `%v` has an undefined rule `%v`", f.Name, r)
	}
	return rulesList[r].Prepare(f, r, p)
}
