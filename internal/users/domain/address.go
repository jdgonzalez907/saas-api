package domain

import "errors"

var (
	ErrInvalidStreet      = errors.New("street address cannot be empty")
	ErrInvalidPostalCode  = errors.New("postal code cannot be empty")
	ErrInvalidCity        = errors.New("city name cannot be empty")
	ErrInvalidState       = errors.New("state name cannot be empty")
	ErrInvalidCountry     = errors.New("country name cannot be empty")
	ErrInvalidDescription = errors.New("address description cannot be empty")
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

func (a Address) Equals(other Address) bool {
	if a.street != other.street || a.city != other.city || a.state != other.state || a.country != other.country {
		return false
	}
	if (a.postalCode == nil) != (other.postalCode == nil) {
		return false
	}
	if a.postalCode != nil && *a.postalCode != *other.postalCode {
		return false
	}
	if (a.description == nil) != (other.description == nil) {
		return false
	}
	if a.description != nil && *a.description != *other.description {
		return false
	}
	return true
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
