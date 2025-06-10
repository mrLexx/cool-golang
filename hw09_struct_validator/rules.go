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
	Type     []reflect.Kind
	validate func(p string, v reflect.Value) error
}

type ListRules map[string]ItemRule

type valueSet struct {
	val reflect.Value
	t   reflect.Kind
}

func validateLen(p string, v reflect.Value) error {
	l, err := strconv.Atoi(p)
	if err != nil {
		return NewExecuteError(ErrExecuteCompileRule, "rule `len` must be int, but `len:%v`", p)
	}
	kn := v.Kind()
	switch {
	case kn == reflect.String:
		if utf8.RuneCountInString(v.Interface().(string)) > l {
			return fmt.Errorf("should be less %v: %w", l, ErrValidationLen)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be string, but (%v) is obtained.", kn)
	}

	return nil
}

func validateRegexp(p string, v reflect.Value) error {
	kn := v.Kind()
	switch {
	case kn == reflect.String:
		if _, ok := regexpList[p]; !ok {
			rg, err := regexp.Compile(p)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"error compile regexp `%v`", p)
			}
			regexpList[p] = rg
		}
		if !regexpList[p].MatchString(v.Interface().(string)) {
			return fmt.Errorf("regexp `%v` not match `%v`: %w", p, v.Interface().(string), ErrValidationRegexp)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be string, but (%v) is obtained.", kn)
	}

	return nil
}

func validateIn(p string, v reflect.Value) error {
	slS := strings.Split(p, ",")
	kn := v.Kind()
	switch {
	case kn == reflect.Int:
		slI := make([]int, len(slS))
		for i, v := range slS {
			v, err := strconv.Atoi(v)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", slS[i])
			}
			slI[i] = v
		}
		if !slices.Contains(slI, v.Interface().(int)) {
			return fmt.Errorf("%v not in %v: %w", v.Interface().(int), slI, ErrValidationIn)
		}
	case kn == reflect.String:
		if !slices.Contains(slS, v.Interface().(string)) {
			return fmt.Errorf("%v not in %v : %w", v.Interface().(string), slS, ErrValidationIn)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int|string, but (%v) is obtained.", kn)
	}
	return nil
}

func validateOut(p string, v reflect.Value) error {
	slS := strings.Split(p, ",")
	kn := v.Kind()
	switch {
	case kn == reflect.Int:
		slI := make([]int, len(slS))
		for i, v := range slS {
			v, err := strconv.Atoi(v)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", slS[i])
			}
			slI[i] = v
		}
		if slices.Contains(slI, v.Interface().(int)) {
			return fmt.Errorf("%v not out %v: %w", v.Interface().(int), slI, ErrValidationIn)
		}
	case kn == reflect.String:
		if slices.Contains(slS, v.Interface().(string)) {
			return fmt.Errorf("%v not out %v : %w", v.Interface().(string), slS, ErrValidationIn)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int|string, but (%v) is obtained.", kn)
	}
	return nil
}

func validateMin(p string, v reflect.Value) error {
	kn := v.Kind()
	switch {
	case kn == reflect.Int:
		m, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		if v.Interface().(int) < m {
			return fmt.Errorf("min %v, but %v : %w", m, v.Interface().(int), ErrValidationMin)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int, but (%v) is obtained.", kn)
	}
	return nil
}

func validateMax(p string, v reflect.Value) error {
	kn := v.Kind()
	switch {
	case kn == reflect.Int:
		m, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		if v.Interface().(int) > m {
			return fmt.Errorf("max %v, but %v : %w", m, v.Interface().(int), ErrValidationMax)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int, but (%v) is obtained.", kn)
	}
	return nil
}
