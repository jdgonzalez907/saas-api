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
	Street      string  `json:"street"`
	PostalCode  *string `json:"postal_code"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Country     string  `json:"country"`
	Description *string `json:"description"`
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
		Street:      street,
		PostalCode:  postalCode,
		City:        city,
		State:       state,
		Country:     country,
		Description: description,
	}, nil
}

func (a Address) String() string {
	str := a.Street + ", " + a.City + ", " + a.State + ", " + a.Country
	if a.Description != nil {
		str += " " + *a.Description
	}
	if a.PostalCode != nil {
		str += " " + *a.PostalCode
	}
	return str
}
