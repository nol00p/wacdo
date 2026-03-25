// Package utils provides shared utility functions such as password validation.
package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
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

// GenerateTempPassword creates a cryptographically random password of the given length
// that satisfies the password validation rules (upper, lower, digit, special).
func GenerateTempPassword(length int) (string, error) {
	const (
		upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lower   = "abcdefghijklmnopqrstuvwxyz"
		digits  = "0123456789"
		special = "!@#$%^&*."
		all     = upper + lower + digits + special
	)

	if length < 4 {
		length = 8
	}

	// Guarantee at least one of each required type
	password := make([]byte, length)
	charsets := []string{upper, lower, digits, special}
	for i, cs := range charsets {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(cs))))
		if err != nil {
			return "", err
		}
		password[i] = cs[idx.Int64()]
	}

	// Fill the rest with random characters from the full set
	for i := 4; i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		password[i] = all[idx.Int64()]
	}

	// Shuffle to avoid predictable positions
	for i := length - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}
