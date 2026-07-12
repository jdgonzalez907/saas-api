package domain

import "errors"

var (
	ErrInvalidStreet      = errors.New("invalid street")
	ErrInvalidPostalCode  = errors.New("invalid postal code")
	ErrInvalidCity        = errors.New("invalid city")
	ErrInvalidState       = errors.New("invalid state")
	ErrInvalidCountry     = errors.New("invalid country")
	ErrInvalidDescription = errors.New("invalid description")
)

type Address struct {
	street      string
	postalCode  *string
	city        string
	state       string
	country     string
	description *string
}

func NewAddress(street, city, state, country string, postalCode, description *string) (Address, error) {
	if postalCode != nil && *postalCode == "" {
		return Address{}, ErrInvalidPostalCode
	}
	if description != nil && *description == "" {
		return Address{}, ErrInvalidDescription
	}
	if street == "" {
		return Address{}, ErrInvalidStreet
	}
	if city == "" {
		return Address{}, ErrInvalidCity
	}
	if state == "" {
		return Address{}, ErrInvalidState
	}
	if country == "" {
		return Address{}, ErrInvalidCountry
	}

	return Address{
		street:      street,
		postalCode:  postalCode,
		city:        city,
		state:       state,
		country:     country,
		description: description,
	}, nil
}

type AddressDTO struct {
	Street      string  `json:"street"`
	PostalCode  *string `json:"postal_code"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Country     string  `json:"country"`
	Description *string `json:"description"`
}

func (a Address) ToDTO() AddressDTO {
	return AddressDTO{
		Street:      a.street,
		PostalCode:  a.postalCode,
		City:        a.city,
		State:       a.state,
		Country:     a.country,
		Description: a.description,
	}
}
