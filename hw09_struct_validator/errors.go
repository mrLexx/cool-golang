package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
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

func NewExecuteError(err error, format string, a ...any) error {
	return &ExecuteError{
		Msg: fmt.Sprintf(format, a...),
		Err: err,
	}
}

var (
	ErrValidationLen    = errors.New("error length string")
	ErrValidationIn     = errors.New("error In")
	ErrValidationOut    = errors.New("error Out")
	ErrValidationMin    = errors.New("error Min")
	ErrValidationMax    = errors.New("error Max")
	ErrValidationRegexp = errors.New("error Regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

func (r *ValidationError) Error() string {
	return fmt.Sprintf("%v: %v", r.Field, r.Err)
}

func (r *ValidationError) Unwrap() error {
	return r.Err
}

func NewValidationError(err error, f string) error {
	return &ValidationError{
		Field: f,
		Err:   err,
	}
}

/*
	 func separateValidateError(err error) error {
		var validErr *ValidationError
		switch {
		case errors.As(err, &validErr):
			validErrs = append(validErrs, *validErr)
		default:
			if err != nil {
				return err
			}
		}
		return nil
	}
*/
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range v {
		sb.WriteString(err.Error())
		if i < len(v)-1 {
			sb.WriteString("; ")
		}
	}
	return sb.String()
}
