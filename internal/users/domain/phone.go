package domain

import (
	"errors"
)

const (
	countryCodeMinLength = 1
	countryCodeMaxLength = 3
	numberMinLength      = 10
	numberMaxLength      = 15
	digitMin             = '0'
	digitMax             = '9'
)

var (
	ErrPhoneEmptyCountryCode    = errors.New("country code cannot be empty")
	ErrPhoneInvalidCountryCode  = errors.New("country code must be between 1 and 3 digits")
	ErrPhoneEmptyNumber         = errors.New("number cannot be empty")
	ErrPhoneInvalidNumber       = errors.New("number must be between 10 and 15 digits")
	ErrPhoneNonDigitCountryCode = errors.New("country code must contain only digits")
	ErrPhoneNonDigitNumber      = errors.New("number must contain only digits")
)

type Phone struct {
	countryCode string
	number      string
}

type PhoneDTO struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

func NewPhone(countryCode, number string) (Phone, error) {
	if countryCode == "" {
		return Phone{}, ErrPhoneEmptyCountryCode
	}

	if len(countryCode) < countryCodeMinLength || len(countryCode) > countryCodeMaxLength {
		return Phone{}, ErrPhoneInvalidCountryCode
	}

	for i := 0; i < len(countryCode); i++ {
		if countryCode[i] < digitMin || countryCode[i] > digitMax {
			return Phone{}, ErrPhoneNonDigitCountryCode
		}
	}

	if number == "" {
		return Phone{}, ErrPhoneEmptyNumber
	}

	if len(number) < numberMinLength || len(number) > numberMaxLength {
		return Phone{}, ErrPhoneInvalidNumber
	}

	for i := 0; i < len(number); i++ {
		if number[i] < digitMin || number[i] > digitMax {
			return Phone{}, ErrPhoneNonDigitNumber
		}
	}

	return Phone{
		countryCode: countryCode,
		number:      number,
	}, nil
}

func (v Phone) CountryCode() string { return v.countryCode }
func (v Phone) Number() string      { return v.number }

func (v Phone) Equals(other Phone) bool {
	return v.countryCode == other.countryCode && v.number == other.number
}

func (v Phone) ToDTO() PhoneDTO {
	return PhoneDTO{
		CountryCode: v.countryCode,
		Number:      v.number,
	}
}
