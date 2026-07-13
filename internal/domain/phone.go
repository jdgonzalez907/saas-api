package domain

import "errors"

var ErrInvalidPhone = errors.New("invalid phone")

type Phone struct {
	countryCode string
	number      string
}

func NewPhone(countryCode, number string) (Phone, error) {
	if countryCode == "" || number == "" {
		return Phone{}, ErrInvalidPhone
	}

	return Phone{
		countryCode: countryCode,
		number:      number,
	}, nil
}

type PhoneDTO struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

func (p Phone) ToDTO() PhoneDTO {
	return PhoneDTO{
		CountryCode: p.countryCode,
		Number:      p.number,
	}
}

func (p Phone) CountryCode() string {
	return p.countryCode
}

func (p Phone) Number() string {
	return p.number
}
