package domain

import (
	"errors"
	"net/mail"
)

var (
	ErrEmailEmpty   = errors.New("email cannot be empty")
	ErrEmailInvalid = errors.New("email is invalid")
)

type Email string

func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", ErrEmailEmpty
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrEmailInvalid
	}

	return Email(email), nil
}

func (v Email) Equals(other Email) bool {
	return v == other
}

func (v Email) ToDTO() string {
	return string(v)
}
