package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	const bSlash = "\\"
	const dblBSlash = "\\\\"

	checkBSlash := func(s string) string {
		if s == dblBSlash {
			return bSlash
		}
		return s
	}

	var unpack, buff, cur, raw strings.Builder

	for _, v := range s {
		cur.Reset()
		cur.WriteRune(v)
		cnt := -1
		isDig := false

		// checking the current character, digit, or not

		if i, err := strconv.Atoi(cur.String()); err == nil {
			isDig, cnt = true, i
		}

		// checking for errors

		if isDig && buff.String() == "" {
			p, c := raw.String(), cur.String()
			return "", fmt.Errorf("digit without symbol or number `%s%s`: %w", p, c, ErrInvalidString)
		}

		if buff.String() == bSlash && !isDig && cur.String() != bSlash {
			p, c := raw.String(), cur.String()
			return "", fmt.Errorf("you can only escape a digit or a slash `%s%s`: %w", p, c, ErrInvalidString)
		}

		// checking the backslash conditions, filling the buffer

		switch {
		case buff.String() == bSlash && isDig: // \n -> "n"
			isDig = false
			buff.Reset()
			buff.WriteString(cur.String())
			cur.Reset()

		case buff.String() == bSlash && cur.String() == bSlash: // \\ -> "\"
			buff.Reset()
			buff.WriteString(dblBSlash)
			cur.Reset()

		case buff.String() == "": // save char in buff
			buff.Reset()
			buff.WriteString(cur.String())
			cur.Reset()
		}

		// unpack string

		switch {
		case isDig:
			unpack.WriteString(strings.Repeat(checkBSlash(buff.String()), cnt))
			buff.Reset()
			isDig = false

		case buff.String() != "" && cur.String() != "":
			unpack.WriteString(checkBSlash(buff.String()))
			buff.Reset()
			buff.WriteString(cur.String())
			cur.Reset()
		}

		raw.Reset()
		raw.WriteRune(v)
	}

	// checking for errors

	if buff.String() == bSlash {
		return "", fmt.Errorf("the backslash is the last one: %w", ErrInvalidString)
	}

	// unpack string

	if buff.String() != "" {
		unpack.WriteString(buff.String())
	}

	return unpack.String(), nil
}
