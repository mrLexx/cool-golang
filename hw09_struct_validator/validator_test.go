package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string          `json:"id" validate:"len:36"`
		Name   string          `json:"name"`
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			test, e := tt.in, tt.expectedErr
			// t.Parallel()

			err := Validate(test)
			fmt.Println(err)
			fmt.Println(e)

			require.Equal(t, err, e)
		})
	}
}

func TestErrorExecute(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			8,
			ErrExecuteWrongInput,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// Undefined
				Somethig []string `validate:"len:18|unknown:18"`
			}{},
			ErrExecuteUndefinedRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// wrong type
				Somethig []string `validate:"min:18"`
			}{},
			ErrExecuteWrongRuleType,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// incomplete
				Somethig []string `validate:"len:|in:as,sd"`
			}{},
			ErrExecuteIncompleteRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// compile
				Somethig []string `validate:"regexp:[a-z"`
			}{},
			ErrExecuteCompileRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// compile
				Somethig []string `validate:"len:df"`
			}{},
			ErrExecuteCompileRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// compile
				Somethig int `validate:"in:as,1,3"`
			}{},
			ErrExecuteCompileRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// compile
				Somethig int `validate:"min:as"`
			}{},
			ErrExecuteCompileRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`

				// compile
				Somethig int `validate:"max:as"`
			}{},
			ErrExecuteCompileRule,
		},
		{
			struct {
				Version string `validate:"len:18"`
				Desc    string `validate:"len:24"`
				None    string `validate:"in:as,sd|len:23"`
				Num     int    `validate:"out:1"`

				// compile
				Somethig int `validate:"out:as"`
			}{},
			ErrExecuteCompileRule,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			test, expectedErr := tt.in, tt.expectedErr
			// t.Parallel()

			err := Validate(test)

			// первый вариант - разворачивания ошибки из пользовательского типа
			var execErr *ExecuteError
			if !errors.As(err, &execErr) {
				t.Fatalf("expected ExecuteError, got %v", err)
			}
			require.ErrorIs(t, execErr.Err, expectedErr)

			// второй вариант: работаем напрямую с ошибкой через Unwrap
			require.ErrorIs(t, err, expectedErr)
			// t.Log(err)
		})
	}
}

func TestExecute(t *testing.T) {
	type UserRoleNestedLevel00 struct {
		Name int    `validate:"len:sd"`
		Desc string `validate:"min:12"`
	}

	type UserRoleNestedLevel01 struct {
		Name string                `validate:"len:12"`
		Desc string                `validate:"len:12"`
		Role UserRoleNestedLevel00 `validate:"nested"`
	}

	type User struct {
		Role    UserRoleNestedLevel01 `validate:"nested"`
		Version string                `validate:"len:12"`
		Email   string                `validate:"len:45|regexp:^\\w+@\\w+\\.\\w+$"`
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{},
			ErrExecuteUndefinedRule,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			_ = t
			test, _ := tt.in, tt.expectedErr
			// t.Parallel()

			if err := Validate(test); err != nil {
				t.Log(err)
			}
		})
	}
}
