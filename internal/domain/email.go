package domain

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email")

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

type EmailDTO struct {
	Value string `json:"value"`
}

func (e Email) ToDTO() EmailDTO {
	return EmailDTO{
		Value: e.value,
	}
}
