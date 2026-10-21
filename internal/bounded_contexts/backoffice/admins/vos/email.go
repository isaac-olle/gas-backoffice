package vos

import (
	"fmt"
	"net/mail"
)

// Email is optional: an empty Email is valid (not every admin has one on
// record), but a non-empty value must be a well-formed address.
type Email string

func NewEmail(raw string) (Email, error) {
	if raw == "" {
		return "", nil
	}
	if _, err := mail.ParseAddress(raw); err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}
	return Email(raw), nil
}

func EmailFromPrimitives(raw string) Email {
	return Email(raw)
}

func (e Email) String() string { return string(e) }
