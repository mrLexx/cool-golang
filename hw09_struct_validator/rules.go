package hw09structvalidator

import (
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

func ValidateLen(p string, v any) error {
	len, err := strconv.Atoi(p)

	if err != nil {
		return NewExecuteError(ErrExecuteCompileRule,
			"rule `len` must be int, but `len:%v`", p)
	}

	switch v := v.(type) {
	case []string:
		for _, v := range v {
			if utf8.RuneCountInString(v) > len {
				return NewValidationError(ErrValidationLen, "")
			}
		}
	case string:
		if utf8.RuneCountInString(v) > len {
			return NewValidationError(ErrValidationLen, "")
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (string)")
	}

	return nil

}
func ValidateRegexp(p string, v any) error {
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
				return NewValidationError(ErrValidationOut, "")
			}
		}
	case string:
		if !regexpList[p].MatchString(v) {
			return NewValidationError(ErrValidationOut, "")
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (string)")
	}

	return nil
}
func ValidateIn(p string, v any) error {

	for _, r := range strings.Split(p, ",") {
		switch v := v.(type) {
		case []string:
			if !slices.Contains(v, r) {
				return NewValidationError(ErrValidationIn, "")
			}
		case string:
			if r != v {
				return NewValidationError(ErrValidationIn, "")
			}
		case []int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", r)
			}
			if !slices.Contains(v, i) {
				return NewValidationError(ErrValidationIn, "")
			}
		case int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", r)
			}
			if i != v {
				return NewValidationError(ErrValidationIn, "")
			}
		default:
			return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (string,int)")

		}
	}
	return nil
}
func ValidateMin(p string, v any) error {
	switch v := v.(type) {
	case []int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		for _, v := range v {
			if v < pi {
				return NewValidationError(ErrValidationMin, "")
			}
		}
	case int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		if v < pi {
			return NewValidationError(ErrValidationMin, "")
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (int)")
	}
	return nil
}
func ValidateMax(p string, v any) error {
	switch v := v.(type) {
	case []int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `max` must be int, but `max:%v`", p)
		}
		for _, v := range v {
			if v > pi {
				return NewValidationError(ErrValidationMax, "")
			}
		}
	case int:
		pi, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `max` must be int, but `max:%v`", p)
		}
		if v > pi {
			return NewValidationError(ErrValidationMax, "")
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (int)")
	}
	return nil
}
func ValidateOut(p string, v any) error {

	for _, r := range strings.Split(p, ",") {
		switch v := v.(type) {
		case []string:
			if slices.Contains(v, r) {
				return NewValidationError(ErrValidationOut, "")
			}
		case string:
			if r == v {
				return NewValidationError(ErrValidationOut, "")
			}
		case []int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `out` must be int, but `out:%v`", r)
			}
			if slices.Contains(v, i) {
				return NewValidationError(ErrValidationOut, "")
			}
		case int:
			i, err := strconv.Atoi(r)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `out` must be int, but `out:%v`", r)
			}
			if i == v {
				return NewValidationError(ErrValidationOut, "")
			}
		default:
			return NewExecuteError(ErrExecuteWrongRuleType, "this rule only for (string,int)")

		}
	}
	return nil
}

type RuleSet struct {
	Name    string
	Payload string
}

func extractRule(tag string) (RuleSet, error) {
	tmp := strings.Split(tag, ":")
	if len(tmp) != 2 || tmp[1] == "" {
		return RuleSet{}, NewExecuteError(ErrExecuteIncompleteRule, "has an incomplete rule `%v`", tag)
	}
	r, p := tmp[0], tmp[1]
	return RuleSet{Name: r, Payload: p}, nil
}

func splitTag(tag string) []string {
	return strings.Split(tag, "|")
}
