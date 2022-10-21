package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
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
	var (
		i    int
		user User
	)
	bs := bufio.NewScanner(r)
	for bs.Scan() {
		if err = json.Unmarshal(bs.Bytes(), &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return result, bs.Err()
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, "."+domain) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
