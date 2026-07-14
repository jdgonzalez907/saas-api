package domain

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("email address must follow a valid format example@domain.com")

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, ErrInvalidEmail
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: email}, nil
}

type EmailDTO string

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

func (e Email) ToDTO() EmailDTO {
	return EmailDTO(e.value)
}

func (e Email) Value() string {
	return e.value
}
