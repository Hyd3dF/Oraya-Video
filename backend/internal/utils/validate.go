package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)

func ValidEmail(s string) bool {
	_, err := mail.ParseAddress(strings.TrimSpace(s))
	return err == nil
}

func ValidUsername(s string) bool {
	return usernameRe.MatchString(s)
}

func ValidPassword(s string) bool {
	return len(s) >= 8 && len(s) <= 128
}
