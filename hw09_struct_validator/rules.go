package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func isInt(v string) bool {
	if _, err := strconv.Atoi(v); err != nil {
		return false
	}
	return true
}

const ValidateTag = "validate"

type Rule struct {
	Type     []reflect.Kind
	Validate func(r, p string, tp reflect.Kind) error
}

type ListRules map[string]Rule

var rulesStore = ListRules{
	"len": {
		Type: []reflect.Kind{reflect.String},
		Validate: func(r, p string, tp reflect.Kind) error {
			_ = tp
			if !isInt(p) {
				return makeExecuteErrorf(ErrExecuteCompileRule,
					"rule `len` must be int, but `%v:%v`", r, p)
			}
			return nil
		},
	},
	"regexp": {
		Type: []reflect.Kind{reflect.String},
		Validate: func(r, p string, tp reflect.Kind) error {
			_ = tp
			if _, ok := regexpList[p]; !ok {
				rg, err := regexp.Compile(p)
				if err != nil {
					return makeExecuteErrorf(ErrExecuteCompileRule,
						"error compile regexp `%v:%v`", r, p)
				}
				regexpList[p] = rg
			}
			return nil
		},
	},
	"in": {
		Type: []reflect.Kind{reflect.String, reflect.Int},
		Validate: func(r, p string, tp reflect.Kind) error {
			if tp == reflect.Int {
				for _, v := range strings.Split(p, ",") {
					if !isInt(v) {
						return makeExecuteErrorf(ErrExecuteCompileRule,
							"rule `in` must be int, but `%v:%v`", r, p)
					}
				}
			}
			return nil
		},
	},
	"min": {
		Type: []reflect.Kind{reflect.Int},
		Validate: func(r, p string, tp reflect.Kind) error {
			_ = tp
			if !isInt(p) {
				return makeExecuteErrorf(ErrExecuteCompileRule,
					"rule `min` must be int, but `%v:%v`", r, p)
			}
			return nil
		},
	},
	"max": {
		Type: []reflect.Kind{reflect.Int},
		Validate: func(r, p string, tp reflect.Kind) error {
			_ = tp
			if !isInt(p) {
				return makeExecuteErrorf(ErrExecuteCompileRule,
					"rule `max` must be int, but `%v:%v`", r, p)
			}
			return nil
		},
	},
	"out": {
		Type: []reflect.Kind{reflect.String, reflect.Int},
		Validate: func(r, p string, tp reflect.Kind) error {
			if tp == reflect.Int {
				for _, v := range strings.Split(p, ",") {
					if !isInt(v) {
						return makeExecuteErrorf(ErrExecuteCompileRule,
							"rule `out` must be int, but `%v:%v`", r, p)
					}
				}
			}
			return nil
		},
	},
}

func extractRule(tag string) (r, p string, err error) {
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
