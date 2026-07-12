package domain

import "errors"

var ErrInvalidPhone = errors.New("invalid phone")

type Phone struct {
	Value string `json:"value"`
}

func NewPhone(phone string) (Phone, error) {
	if phone == "" {
		return Phone{}, ErrInvalidPhone
	}

	return Phone{Value: phone}, nil
}

func (p Phone) String() string {
	return p.Value
}
