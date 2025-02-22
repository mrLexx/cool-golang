package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

const bSlash = "\\"

const dblBSlash = "\\\\"

func Unpack(s string) (string, error) {
	checkBSlash := func(s string) string {
		if s == dblBSlash {
			return bSlash
		}
		return s
	}

	var unpack strings.Builder
	var buff, raw string

	for _, v := range s {
		cur := string(v)

		isDig, cnt := false, -1

		// checking the current character, digit, or not

		if i, err := strconv.Atoi(cur); err == nil {
			isDig, cnt = true, i
		}

		// checking for errors

		if isDig && buff == "" {
			p, c := raw, cur
			return "", fmt.Errorf("digit without symbol or number `%s%s`: %w", p, c, ErrInvalidString)
		}

		if buff == bSlash && !isDig && cur != bSlash {
			p, c := raw, cur
			return "", fmt.Errorf("you can only escape a digit or a slash `%s%s`: %w", p, c, ErrInvalidString)
		}

		// checking the backslash conditions, filling the buffer

		raw = cur

		switch {
		case buff == bSlash && isDig: // \n -> "n"
			isDig = false
			buff = cur
			cur = ""
		case buff == bSlash && cur == bSlash: // \\ -> "\"
			buff = dblBSlash
			cur = ""
		case buff == "": // save char in buff
			buff = cur
			cur = ""
		}

		// unpack string

		switch {
		case isDig:
			unpack.WriteString(strings.Repeat(checkBSlash(buff), cnt))
			buff = ""

		case buff != "" && cur != "":
			unpack.WriteString(checkBSlash(buff))
			buff = cur
		}
	}

	// checking for errors

	if buff == bSlash {
		return "", fmt.Errorf("the backslash is the last one: %w", ErrInvalidString)
	}

	// unpack string

	if buff != "" {
		unpack.WriteString(buff)
	}

	return unpack.String(), nil
}
