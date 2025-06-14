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

type ItemRule struct {
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
		if utf8.RuneCountInString(v.Val.String()) != l {
			return fmt.Errorf("length must by equil %v: %w", l, ErrValidationLen)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (string), but (%v) is obtained.", v.Type)
	}

	return nil
}

func validateRegexp(p string, v valueSet) error {
	switch {
	case v.Type == reflect.String:
		rg, ok := regexpCache.Get(cacheValue{Rule: p})
		if !ok {
			trg, err := regexp.Compile(p)
			if err != nil {
				return NewExecuteError(ErrExecuteCompileRule,
					"error compile regexp `%v`", p)
			}
			regexpCache.Set(cacheValue{Rule: p}, trg)
			rg = trg
		}

		if !rg.MatchString(v.Val.String()) {
			return fmt.Errorf("regexp `%v` not match `%v`: %w", p, v.Val.String(), ErrValidationRegexp)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (string), but (%v) is obtained.", v.Type)
	}

	return nil
}

func validateIn(p string, v valueSet) error {
	switch {
	case v.Type == reflect.Int:
		if vc, ok := inCache.Get(cacheValue{Val: fmt.Sprint(v.Val.Int()), Rule: p}); ok {
			if vc {
				return nil
			}
			return fmt.Errorf("%v not in %v: %w", v.Val.Int(), p, ErrValidationIn)
		}

		slS := strings.Split(p, ",")
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
			inCache.Set(cacheValue{Val: fmt.Sprint(v.Val.Int()), Rule: p}, false)
			return fmt.Errorf("%v not in %v: %w", v.Val.Int(), p, ErrValidationIn)
		}
		inCache.Set(cacheValue{Val: fmt.Sprint(v.Val.Int()), Rule: p}, true)
	case v.Type == reflect.String:
		if vc, ok := inCache.Get(cacheValue{Val: v.Val.String(), Rule: p}); ok {
			if vc {
				return nil
			}
			return fmt.Errorf("%v not in %v: %w", v.Val.String(), p, ErrValidationIn)
		}

		slS := strings.Split(p, ",")
		if !slices.Contains(slS, v.Val.String()) {
			inCache.Set(cacheValue{Val: v.Val.String(), Rule: p}, false)
			return fmt.Errorf("%v not in %v : %w", v.Val.String(), slS, ErrValidationIn)
		}
		inCache.Set(cacheValue{Val: v.Val.String(), Rule: p}, true)
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (int,string), but (%v) is obtained.", v.Type)
	}
	return nil
}

func validateOut(p string, v valueSet) error {
	switch {
	case v.Type == reflect.Int:
		if vc, ok := outCache.Get(cacheValue{Val: fmt.Sprint(v.Val.Int()), Rule: p}); ok {
			if !vc {
				return nil
			}
			return fmt.Errorf("%v not out %v: %w", v.Val.Int(), p, ErrValidationOut)
		}

		slS := strings.Split(p, ",")
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
			outCache.Set(cacheValue{Val: fmt.Sprint(v.Val.Int()), Rule: p}, true)
			return fmt.Errorf("%v not out %v: %w", v.Val.Int(), slI, ErrValidationOut)
		}
		outCache.Set(cacheValue{Val: fmt.Sprint(v.Val.Int()), Rule: p}, false)
	case v.Type == reflect.String:
		if vc, ok := outCache.Get(cacheValue{Val: v.Val.String(), Rule: p}); ok {
			if !vc {
				return nil
			}
			return fmt.Errorf("%v not out %v: %w", v.Val.String(), p, ErrValidationOut)
		}

		slS := strings.Split(p, ",")
		if slices.Contains(slS, v.Val.String()) {
			outCache.Set(cacheValue{Val: v.Val.String(), Rule: p}, true)
			return fmt.Errorf("%v not out %v : %w", v.Val.String(), slS, ErrValidationOut)
		}
		outCache.Set(cacheValue{Val: v.Val.String(), Rule: p}, false)
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (int,string), but (%v) is obtained.", v.Type)
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
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (int), but (%v) is obtained.", v.Type)
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
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (int), but (%v) is obtained.", v.Type)
	}
	return nil
}

func validateRequire(p string, v valueSet) error {
	r, err := strconv.ParseBool(p)
	if err != nil {
		return NewExecuteError(ErrExecuteCompileRule,
			"rule `require` must be bool, but `require:%v`", p)
	}

	switch {
	case v.Type == reflect.Int:
		if r && v.Val.IsZero() {
			return fmt.Errorf("is zero %v : %w", v.Val.Int(), ErrValidationRequire)
		}
	case v.Type == reflect.String:
		if r && v.Val.IsZero() {
			return fmt.Errorf("is zero %v : %w", v.Val.String(), ErrValidationRequire)
		}
	default:
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be (int,string), but (%v) is obtained.", v.Type)
	}
	return nil
}
