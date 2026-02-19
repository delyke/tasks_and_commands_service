package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")

	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// Email represents a validated email address value object.
type Email string

// NewEmail creates and validates an Email value object.
func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if !emailRegex.MatchString(value) {
		return "", ErrInvalidEmail
	}
	return Email(value), nil
}

// String returns the string representation of the email.
func (e Email) String() string {
	return string(e)
}
