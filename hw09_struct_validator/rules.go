package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	Type    []reflect.Kind
	Prepare func(f reflect.StructField, r, p string) error
}

func isInt(v string) bool {
	if _, err := strconv.Atoi(v); err != nil {
		return false
	}
	return true
}

const ValidateTag = "validate"

var rulesList = map[string]Rule{
	"len": {
		Type: []reflect.Kind{reflect.String},
		Prepare: func(f reflect.StructField, r, p string) error {
			if !isInt(p) {
				return makeExecuteErrorf(ErrExecuteCompileRule,
					"field `%v`: tag `len` must be int, but `%v:%v`", f.Name, r, p)
			}
			return nil
		},
	},
	"regexp": {
		Type: []reflect.Kind{reflect.String},
		Prepare: func(f reflect.StructField, r, p string) error {
			_ = r
			if _, ok := regexpList[p]; !ok {
				rg, err := regexp.Compile(p)
				if err != nil {
					return makeExecuteErrorf(ErrExecuteCompileRule,
						"field `%v` error compile regexp `%v`", f.Name, p)
				}
				regexpList[p] = rg
			}
			return nil
		},
	},
	"in": {
		Type: []reflect.Kind{reflect.String, reflect.Int},
		Prepare: func(f reflect.StructField, r, p string) error {
			if getReflectType(f) == reflect.Int {
				for _, v := range strings.Split(p, ",") {
					if !isInt(v) {
						return makeExecuteErrorf(ErrExecuteCompileRule,
							"field `%v`: tag `in` must be int, but `%v:%v`", f.Name, r, p)
					}
				}
			}
			return nil
		},
	},
	"min": {
		Type: []reflect.Kind{reflect.Int},
		Prepare: func(f reflect.StructField, r, p string) error {
			if !isInt(p) {
				return makeExecuteErrorf(ErrExecuteCompileRule,
					"field `%v`: tag `min` must be int, but `%v:%v`", f.Name, r, p)
			}
			return nil
		},
	},
	"max": {
		Type: []reflect.Kind{reflect.Int},
		Prepare: func(f reflect.StructField, r, p string) error {
			if !isInt(p) {
				return makeExecuteErrorf(ErrExecuteCompileRule,
					"field `%v`: tag `max` must be int, but `%v:%v`", f.Name, r, p)
			}
			return nil
		},
	},
	"out": {
		Type: []reflect.Kind{reflect.String, reflect.Int},
		Prepare: func(f reflect.StructField, r, p string) error {
			if getReflectType(f) == reflect.Int {
				for _, v := range strings.Split(p, ",") {
					if !isInt(v) {
						return makeExecuteErrorf(ErrExecuteCompileRule,
							"field `%v`: tag `out` must be int, but `%v:%v`", f.Name, r, p)
					}
				}
			}
			return nil
		},
	},
}

var regexpList = make(map[string]*regexp.Regexp, 0)
