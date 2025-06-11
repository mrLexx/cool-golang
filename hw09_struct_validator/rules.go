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
	Typies   []reflect.Kind
	Validate func(p string, v valueSet, allowedTp []reflect.Kind) error
}

type ListRules map[string]ItemRule

type valueSet struct {
	Val  reflect.Value
	Type reflect.Kind
}

func checkingTypes(allowedTp []reflect.Kind, tp reflect.Kind) error {
	if !slices.Contains(allowedTp, tp) {
		return NewExecuteError(ErrExecuteWrongRuleType, "the type must be %v, but (%v) is obtained.", allowedTp, tp)
	}
	return nil
}

func validateLen(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}

	return nil
}

func validateRegexp(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}

	return nil
}

func validateIn(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}
	return nil
}

func validateOut(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
			return fmt.Errorf("%v not out %v: %w", v.Val.Int(), slI, ErrValidationOut)
		}
	case v.Type == reflect.String:
		if slices.Contains(slS, v.Val.String()) {
			return fmt.Errorf("%v not out %v : %w", v.Val.String(), slS, ErrValidationOut)
		}
	default:
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}
	return nil
}

func validateMin(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}
	return nil
}

func validateMax(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}
	return nil
}

func validateRequire(p string, v valueSet, allowedTp []reflect.Kind) error {
	if err := checkingTypes(allowedTp, v.Type); err != nil {
		return err
	}

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
		panic(fmt.Sprintf("release validate type from %v", allowedTp))
	}
	return nil
}
