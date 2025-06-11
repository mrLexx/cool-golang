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
	MetaSub struct {
		Ident []string `validate:"len:5"`
		Desc  string   `validate:"require"`
		Point int      `validate:"require"`
	}
	Meta struct {
		Info  string    `validate:"len:17"`
		Range int       `validate:"min:10|max:50|out:45,27"`
		Sub   []MetaSub `validate:"nested"`
	}
	User struct {
		ID     string `json:"id" validate:"len:15"`
		Name   string
		Age    int      `validate:"min:18|max:50|out:23"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff|out:root"`
		Phones []string `validate:"len:10"`
		Limbs  int      `validate:"in:1,2,3,4"`
		Eyes   int      `validate:"min:0|max:4|out:3,4"`
		Weight int      `validate:"require"`
		Bio    string   `validate:"require"`
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
)

//nolint:funlen
func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			Token{
				Header:    []byte("header data"),
				Payload:   []byte("payload data"),
				Signature: []byte("signature data"),
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "sdsd",
			},
			nil,
		},
		{
			Response{
				Code: 202,
				Body: "sdsd",
			},
			ErrValidationIn,
		},
		{
			[]Response{
				{
					Code: 200,
					Body: "sdsd",
				},
				{
					Code: 202,
					Body: "sdsd",
				},
			},
			ErrValidationIn,
		},
		{
			App{
				Version: "12345",
			},
			nil,
		},
		{
			App{
				Version: "123456",
			},
			ErrValidationLen,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			nil,
		},
		{
			User{
				ID:     "HASH😎678901234", // error
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationLen,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "emailmail.ru", // error
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationRegexp,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "other", // error
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationIn,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "root",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationOut,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationLen,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  6, // error
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationIn,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 0, // error
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationRequire,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "", // error
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationRequire,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 45, // error
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationOut,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 51, // error
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationMax,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 9, // error
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationMin,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "nformation abou👍", // error
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationLen,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"2345", // error
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcv", // error
							},
							Desc:  "go, golang, goshechka",
							Point: 2,
						},
					},
				},
			},
			ErrValidationLen,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							// Desc:  "go, golang, goshechka",// error
							Point: 2,
						},
					},
				},
			},
			ErrValidationRequire,
		},
		{
			User{
				ID:     "HASH😎6789012345",
				Name:   "User name 🙌",
				Age:    18,
				Email:  "email@mail.ru",
				Role:   "admin",
				Phones: []string{"9652025404", "9601044485"},
				Limbs:  4,
				Eyes:   2,
				Weight: 85,
				Bio:    "Bio f f fg fg f fg",
				Meta: Meta{
					Info:  "Information abou👍",
					Range: 25,
					Sub: []MetaSub{
						{
							Ident: []string{
								"12345",
								"67890",
							},
							Desc:  "Description about this thing",
							Point: 1,
						},
						{
							Ident: []string{
								"asdfg",
								"zxcvb",
							},
							Desc: "go, golang, goshechka",
							// Point: 2,// error
						},
					},
				},
			},
			ErrValidationRequire,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			_ = t
			test, expectedErr := tt.in, tt.expectedErr
			_ = expectedErr
			t.Parallel()

			err := Validate(test)

			switch {
			case expectedErr == nil:
				require.Nil(t, err)
			default:
				var validErrs ValidationErrors
				if !errors.As(err, &validErrs) {
					t.Fatalf("expected ValidationErrors, got %v", err)
				}
				require.ErrorIs(t, validErrs, expectedErr)
			}
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
