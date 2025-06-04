package hw09structvalidator

import (
	"errors"
	"fmt"
)

var (
	ErrExecuteIncompleteRule = errors.New("error incomplete rule")
	ErrExecuteUndefinedRule  = errors.New("error undefined rule")
	ErrExecuteWrongInput     = errors.New("error wrong input")
	ErrExecuteWrongRuleType  = errors.New("error wrong rule type")
	ErrExecuteCompileRule    = errors.New("error compile rule")
)

type ExecuteError struct {
	Msg string
	Err error
}

func (r *ExecuteError) Error() string {
	return fmt.Sprintf("%v: %v", r.Msg, r.Err)
}

func (r *ExecuteError) Unwrap() error {
	return r.Err
}

func makeExecuteErrorf(err error, format string, a ...any) error {
	return &ExecuteError{
		Msg: fmt.Sprintf(format, a...),
		Err: err,
	}
}

type ValidationError struct {
	Field string
	Err   error
}

func (r *ValidationError) Error() string {
	return ""
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return "many many validate errors occurred"
}
