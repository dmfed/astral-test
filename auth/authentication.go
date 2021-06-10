package auth

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

var userPasswd = regexp.MustCompile(`(\w+) (.+)`)

var ErrNoCredentials = errors.New("no credentials found in file")

type Authenticator interface {
	CredentialsAreValid(username, password string) bool
}

type CredentialsKeeper map[string]string

func New(filename string) (Authenticator, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return newCK(file)
}

func (ck *CredentialsKeeper) CredentialsAreValid(username, password string) bool {
	if pwd, ok := (*ck)[username]; ok && password == pwd {
		return true
	}
	return false
}

func newCK(r io.Reader) (Authenticator, error) {
	var ck CredentialsKeeper
	b, err := io.ReadAll(r)
	if err != nil {
		return &ck, err
	}
	userMap := stringsToMap(bytesToStringSlice(b))
	if len(userMap) == 0 {
		return &ck, ErrNoCredentials
	}
	ck = userMap
	return &ck, nil
}

func bytesToStringSlice(b []byte) []string {
	return strings.Split(string(b), "\n")
}

func stringsToMap(s []string) map[string]string {
	m := make(map[string]string)
	for _, line := range s {
		if userPasswd.MatchString(line) {
			username, password := userPasswd.FindStringSubmatch(line)[1], userPasswd.FindStringSubmatch(line)[2]
			m[username] = password
		}
	}
	return m
}
