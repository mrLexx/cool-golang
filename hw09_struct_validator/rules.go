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
	Validate func(p string, v valueSet) error
}

type ListRules map[string]ItemRule

type valueSet struct {
	Val  reflect.Value
	Type reflect.Kind
}

func validateLen(p string, v valueSet) error {
	l, err := strconv.Atoi(p)
	if err != nil {
		return NewExecuteError(ErrExecuteCompileRule, "rule `len` must be int, but `len:%v`", p)
	}
	switch {
	case v.Type == reflect.String:
		if utf8.RuneCountInString(v.Val.String()) > l {
			return fmt.Errorf("should be less %v: %w", l, ErrValidationLen)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be string, but (%v) is obtained.", v.Type)
	}

	return nil
}

func validateRegexp(p string, v valueSet) error {
	switch {
	case v.Type == reflect.String:
		if _, ok := regexpList[p]; !ok {
			rg, err := regexp.Compile(p)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"error compile regexp `%v`", p)
			}
			regexpList[p] = rg
		}
		if !regexpList[p].MatchString(v.Val.String()) {
			return fmt.Errorf("regexp `%v` not match `%v`: %w", p, v.Val.String(), ErrValidationRegexp)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be string, but (%v) is obtained.", v.Type)
	}

	return nil
}

func validateIn(p string, v valueSet) error {
	slS := strings.Split(p, ",")
	switch {
	case v.Type == reflect.Int:
		slI := make([]int64, len(slS))
		for i, v := range slS {
			v, err := strconv.Atoi(v)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", slS[i])
			}
			slI[i] = int64(v)
		}
		if !slices.Contains(slI, v.Val.Int()) {
			return fmt.Errorf("%v not in %v: %w", v.Val.Int(), slI, ErrValidationIn)
		}
	case v.Type == reflect.String:
		if !slices.Contains(slS, v.Val.String()) {
			return fmt.Errorf("%v not in %v : %w", v.Val.String(), slS, ErrValidationIn)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int|string, but (%v) is obtained.", v.Type)
	}
	return nil
}

func validateOut(p string, v valueSet) error {
	slS := strings.Split(p, ",")
	switch {
	case v.Type == reflect.Int:
		slI := make([]int64, len(slS))
		for i, v := range slS {
			v, err := strconv.Atoi(v)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"rule `in` must be int, but `in:%v`", slS[i])
			}
			slI[i] = int64(v)
		}
		if slices.Contains(slI, v.Val.Int()) {
			return fmt.Errorf("%v not out %v: %w", v.Val.Int(), slI, ErrValidationIn)
		}
	case v.Type == reflect.String:
		if slices.Contains(slS, v.Val.String()) {
			return fmt.Errorf("%v not out %v : %w", v.Val.String(), slS, ErrValidationIn)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int|string, but (%v) is obtained.", v.Type)
	}
	return nil
}

func validateMin(p string, v valueSet) error {
	switch {
	case v.Type == reflect.Int:
		m, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		if v.Val.Int() < int64(m) {
			return fmt.Errorf("min %v, but %v : %w", m, v.Val.Int(), ErrValidationMin)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int, but (%v) is obtained.", v.Type)
	}
	return nil
}

func validateMax(p string, v valueSet) error {
	switch {
	case v.Type == reflect.Int:
		m, err := strconv.Atoi(p)
		if err != nil {
			return NewExecuteError(ErrExecuteCompileRule,
				"rule `min` must be int, but `min:%v`", p)
		}
		if v.Val.Int() > int64(m) {
			return fmt.Errorf("max %v, but %v : %w", m, v.Val.Int(), ErrValidationMax)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be int, but (%v) is obtained.", v.Type)
	}
	return nil
}
