package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

var regexpList = make(map[string]*regexp.Regexp, 0)

type ItemRule struct {
	Type           []reflect.Kind
	ValidateStruct func(p string, tp reflect.Kind) error
	ValidateData   func(p string, v any) error
}

type ListRules map[string]ItemRule

func validateLen(p string, v any) error {
	l, err := strconv.Atoi(p)
	if err != nil {
		return NewExecuteError(ErrExecuteCompileRule,
			"rule `len` must be int, but `len:%v`", p)
	}

	switch v := v.(type) {
	case []string:
		for _, v := range v {
			if utf8.RuneCountInString(v) > l {
				return ErrValidationLen
			}
		}
	case string:
		if utf8.RuneCountInString(v) > l {
			return ErrValidationLen
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be string, but (%T) is obtained.", v)
	}

	return nil
}

func validateRegexp(p string, v any) error {
	if _, ok := regexpList[p]; !ok {
		rg, err := regexp.Compile(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"error compile regexp `%v`", p)
		}
		regexpList[p] = rg
	}

	switch v := v.(type) {
	case []string:
		for _, v := range v {
			if !regexpList[p].MatchString(v) {
				return ErrValidationRegexp
			}
		}
	case string:
		if !regexpList[p].MatchString(v) {
			return ErrValidationRegexp
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be string, but (%T) is obtained.", v)
	}

	return nil
}

func validateIn(p string, v any) error {
	rv := reflect.TypeOf(v)
	fmt.Printf("type: %v\n", rv.Kind())
	if rv.Kind() == reflect.Slice {
		fmt.Printf("\tsub type: %v\n", rv.Elem().Kind())
	}
	fmt.Printf("type: %T\n", v)
	fmt.Printf("\n")

	for _, r := range strings.Split(p, ",") {
		switch v := v.(type) {

		case []string:
			if !slices.Contains(v, r) {
				return ErrValidationIn
			}
		case string:
			if r != v {
				return ErrValidationIn
			}
		case []int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", r)
			}
			if !slices.Contains(v, i) {
				return ErrValidationIn
			}
		case int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", r)
			}
			if i != v {
				return ErrValidationIn
			}
		default:
			fmt.Printf("type: %v\n", reflect.ValueOf(v).Kind())
			fmt.Printf("type: %T\n", v)
			return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int|string, but (%T) is obtained.", v)
		}
	}
	return nil
}

func validateMin(p string, v any) error {
	switch v := v.(type) {
	case []int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		for _, v := range v {
			if v < pi {
				return ErrValidationMin
			}
		}
	case int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		if v < pi {
			return ErrValidationMin
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int, but (%T) is obtained.", v)
	}
	return nil
}

func validateMax(p string, v any) error {
	switch v := v.(type) {
	case []int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `max` must be int, but `max:%v`", p)
		}
		for _, v := range v {
			if v > pi {
				return ErrValidationMax
			}
		}
	case int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `max` must be int, but `max:%v`", p)
		}
		if v > pi {
			return ErrValidationMax
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int, but (%T) is obtained.", v)
	}
	return nil
}

func validateOut(p string, v any) error {
	for _, r := range strings.Split(p, ",") {
		switch v := v.(type) {
		case []string:
			if slices.Contains(v, r) {
				return ErrValidationOut
			}
		case string:
			if r == v {
				return ErrValidationOut
			}
		case []int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `out` must be int, but `out:%v`", r)
			}
			if slices.Contains(v, i) {
				return ErrValidationOut
			}
		case int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `out` must be int, but `out:%v`", r)
			}
			if i == v {
				return ErrValidationOut
			}
		default:
			return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int|string, but (%T) is obtained.", v)
		}
	}
	return nil
}

type ruleSet struct {
	Name    string
	Payload string
}
