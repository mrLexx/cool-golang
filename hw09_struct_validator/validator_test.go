package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type UserRoleNested struct {
	Role UserRole `validate:"in:admin,stuff"`
}

// Test the function on different structures and other types.
type (
	Meta struct {
		Info  string `validate:"len:11"`
		Range int    `validate:"min:10|max:50"`
	}
	User struct {
		ID     string   `json:"id" validate:"len:36|regexp:^\\w+@\\w+\\.\\w+$"`
		Phones []string `validate:"len:11"`
		Name   string   `json:"name"`
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Meta   Meta     `validate:"nested"`
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

	MyNested struct {
		OtherF []string `json:"id" validate:"len:5"`
	}

	My struct {
		InString  []string `validate:"in:200,404,500"`
		InInt     []int    `validate:"in:200,404,500"`
		OutString []string `validate:"out:200,404,500"`
		OutInt    []int    `validate:"out:200,404,500"`
		Email     string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Len       string   `validate:"len:3"`
		// Role []UserRoleNested `validate:"nested"`
		// ID     string   `json:"id" validate:"len:36|regexp:^\\w+@\\w+\\.\\w+$"`

		// ID     string     `json:"id" validate:"len:1"`
		// Phones []string   `validate:"len:15"`
		// Nested []MyNested `validate:"nested"`
	}
)

func TestExecute(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			My{
				InString:  []string{"200", "404"},
				InInt:     []int{200, 404},
				OutString: []string{"201", "300"},
				OutInt:    []int{201, 300},
				Email:     "as@as.com",
				Len:       "1234",

				// Role: []UserRoleNested{
				// 	{Role: "admin"},
				// },
				/* ID: "dd👍",
				Phones: []string{
					"12345678901👍",
					"phont2",
				},

				Nested: []MyNested{
					{
						OtherF: []string{"phont1👍", "phont2"},
					},
					{
						OtherF: []string{"level1👍", "level2"},
					},
				}, */
			}, ErrValidationLen,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			_ = t
			test, expectedErr := tt.in, tt.expectedErr
			_ = expectedErr
			t.Parallel()

			err := Validate(test)

			var validErr *ValidationError
			var execErr *ExecuteError
			switch {
			case errors.As(err, &validErr):
				t.Log("Valid error!")
				t.Log(err)
			case errors.As(err, &execErr):
				t.Log("Execute error!")
				t.Log(err)
			}

			// if !errors.As(err, &validErr) {
			// 	t.Fatalf("expected ValidationErrors, got %v", err)
			// }
			// require.ErrorIs(t, err, expectedErr)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:   "dd",
				Name: "name",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			_ = t
			test, expectedErr := tt.in, tt.expectedErr
			_ = expectedErr
			// t.Parallel()

			err := Validate(test)
			_ = err
			// var validErrs *ValidationError
			// if !errors.As(err, &validErrs) {
			// t.Fatalf("expected ValidationErrors, got %v", err)
			// }
			// require.ErrorIs(t, err, expectedErr)
		})
	}
}

func TestErrorExecute(t *testing.T) {
	type UserRoleNestedLevel00 struct {
		Name int    `validate:"len:12"`
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
		{
			User{},
			ErrExecuteWrongRuleType,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			test, expectedErr := tt.in, tt.expectedErr
			t.Parallel()

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
