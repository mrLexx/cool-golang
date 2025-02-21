package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		{input: "a5", expected: "aaaaa"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: "٠١2٢3٣0٤٥٦٧٨٩", expected: "٠١١٢٢٢٤٥٦٧٨٩"},
		{input: `٠١2٢3٣0٤٥٦٧٨٩`, expected: `٠١١٢٢٢٤٥٦٧٨٩`},
		{input: "р1у2с3ские0 буквы", expected: "руусссски буквы"},
		{input: "иероглифы ト0ヨ2タ自動車株式会社3", expected: "иероглифы ヨヨタ自動車株式会社社社"},
		{input: "🚘トヨタ自動車株式会社 🧨2👍5", expected: "🚘トヨタ自動車株式会社 🧨🧨👍👍👍👍👍"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `qwe\\\`, `qw\ne`}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
