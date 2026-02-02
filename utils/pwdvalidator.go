package utils

import (
	"errors"
	"regexp"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password needs to be at least 8 characters long")
	}

	var (
		upper   = regexp.MustCompile(`[A-Z]`)
		lower   = regexp.MustCompile(`[a-z]`)
		number  = regexp.MustCompile(`[0-9]`)
		special = regexp.MustCompile(`[!@#\$%\^&\*\.]`)
	)

	if !upper.MatchString(password) {
		return errors.New("Password Not Compliant: min 1 Maj")
	}
	if !lower.MatchString(password) {
		return errors.New("Password Not Compliant: min 1 lower case")
	}
	if !number.MatchString(password) {
		return errors.New("Password Not Compliant: min 1 Number")
	}
	if !special.MatchString(password) {
		return errors.New("Password Not Compliant: min 1 special")
	}

	return nil
}
