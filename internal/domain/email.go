package domain

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email")

type Email struct {
	Value string `json:"value"`
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, ErrInvalidEmail
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}

	return Email{Value: email}, nil
}

func (e Email) String() string {
	return e.Value
}
