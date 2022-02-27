package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = "." + domain
	n := len(domain)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	reader := bufio.NewReader(r)
	result := make(DomainStat)
	for {
		line, err1 := reader.ReadBytes('\n')
		if len(line) > 0 {
			var user User
			if err := json.Unmarshal(line, &user); err != nil {
				return nil, err
			}
			offset := len(user.Email) - n
			if offset < 0 {
				continue
			}
			if user.Email[offset:] == domain {
				key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
				result[key]++
			}
		}
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			return nil, err1
		}
	}
	return result, nil
}
