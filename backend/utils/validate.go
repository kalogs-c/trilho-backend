package utils

import (
	"regexp"
	"strings"
)

func ValidateName(name string) bool {
	r, _ := regexp.Compile(`^[A-Z][a-z]{1,}$`)

	return r.MatchString(name)
}

func ValidatePassword(password string) bool {
	return len(password) > 5 && !strings.Contains(password, " ")
}
