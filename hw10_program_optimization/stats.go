package hw10programoptimization

import (
	"bufio"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)

	var (
		user User
		l    string
	)

	d := "." + domain

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l = scanner.Text()
		if strings.Contains(l, d) {
			if err = json.Unmarshal([]byte(l), &user); err != nil {
				return
			}
			if strings.Contains(user.Email, d) {
				result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
			}
		}
	}
	return
}
