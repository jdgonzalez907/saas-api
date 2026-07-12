package domain

import "errors"

var ErrInvalidPhone = errors.New("invalid phone")

type Phone struct {
	value string
}

func NewPhone(phone string) (Phone, error) {
	if phone == "" {
		return Phone{}, ErrInvalidPhone
	}

	return Phone{value: phone}, nil
}

type PhoneDTO struct {
	Value string `json:"value"`
}

func (p Phone) ToDTO() PhoneDTO {
	return PhoneDTO{
		Value: p.value,
	}
}
