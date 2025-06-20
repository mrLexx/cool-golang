package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	//nolint:depguard
	"github.com/goccy/go-json"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

var (
	ErrBrokenEmail = errors.New("broken email")
	ErrBrokenJSON  = errors.New("broken JSON")
	ErrReadJSON    = errors.New("error read JSON")
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	var user User

	result := make(DomainStat)

	d := "." + domain

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), d) {
			continue
		}

		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			var jsonErr *json.SyntaxError
			switch {
			case errors.As(err, &jsonErr):
				problemPart := problemJSON(scanner.Bytes(), jsonErr.Offset)
				return nil, fmt.Errorf("error near '%s' (offset %d): %w", problemPart, jsonErr.Offset, ErrBrokenJSON)
			default:
				return nil, fmt.Errorf("broken json: %w", ErrBrokenJSON)
			}
		}

		if !strings.HasSuffix(user.Email, d) {
			continue
		}

		sp := strings.SplitN(user.Email, "@", 2)
		if len(sp) != 2 {
			return nil, fmt.Errorf("broken email %v: %w", user.Email, ErrBrokenEmail)
		}

		dm := strings.ToLower(sp[1])
		result[dm]++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error read json (%v): %w", err.Error(), ErrReadJSON)
	}

	return result, nil
}

func problemJSON(json []byte, offset int64) []byte {
	from := min(offset, 10)
	to := min(int64(len(json))-offset, 10)
	return json[offset-from : offset+to]
}
