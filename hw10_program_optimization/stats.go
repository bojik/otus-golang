package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	reader := bufio.NewReader(r)
	for i := 0; ; i++ {
		line, err1 := reader.ReadBytes('\n')
		if len(line) > 0 {
			var user User
			if err = json.Unmarshal(line, &user); err != nil {
				return
			}
			result[i] = user
		}
		if err1 == io.EOF {
			return
		}
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	domain = "." + domain
	n := len(domain)
	result := make(DomainStat)
	for _, user := range u {
		s := len(user.Email) - n
		if s < 0 {
			s = 0
		}
		matched := user.Email[s:] == domain
		if matched {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}
	return result, nil
}
