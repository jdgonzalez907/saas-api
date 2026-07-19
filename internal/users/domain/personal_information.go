package domain

import (
	"errors"
	"time"
)

type DNIType string

const (
	DNITypeCC  DNIType = "CC"
	DNITypeCE  DNIType = "CE"
	DNITypeNIT DNIType = "NIT"
	DNITypePP  DNIType = "PP"
)

var (
	ErrPersonalInformationInvalidDNIType   = errors.New("DNI type must be one of: CC, CE, NIT, PP")
	ErrPersonalInformationInvalidDNINumber = errors.New("DNI number must be between 8 and 20 characters")
	ErrPersonalInformationInvalidFirstName = errors.New("first name must be between 4 and 100 characters")
	ErrPersonalInformationInvalidLastName  = errors.New("last name must be between 4 and 100 characters")
	ErrPersonalInformationEmptyBirthdate   = errors.New("birthdate cannot be empty")
	ErrPersonalInformationUnderage         = errors.New("must be 18 years or older")
)

type PersonalInformation struct {
	dniType   DNIType
	dniNumber string
	firstName string
	lastName  string
	birthdate time.Time
}

type PersonalInformationDTO struct {
	DNIType   DNIType `json:"dni_type"`
	DNINumber string  `json:"dni_number"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Birthdate string  `json:"birthdate"`
}

func NewPersonalInformation(
	dniType DNIType,
	dniNumber, firstName, lastName string,
	birthdate time.Time,
) (PersonalInformation, error) {
	switch dniType {
	case DNITypeCC, DNITypeCE, DNITypeNIT, DNITypePP:
	default:
		return PersonalInformation{}, ErrPersonalInformationInvalidDNIType
	}

	if len(dniNumber) < 8 || len(dniNumber) > 20 {
		return PersonalInformation{}, ErrPersonalInformationInvalidDNINumber
	}

	if len(firstName) < 4 || len(firstName) > 100 {
		return PersonalInformation{}, ErrPersonalInformationInvalidFirstName
	}

	if len(lastName) < 4 || len(lastName) > 100 {
		return PersonalInformation{}, ErrPersonalInformationInvalidLastName
	}

	if birthdate.IsZero() {
		return PersonalInformation{}, ErrPersonalInformationEmptyBirthdate
	}

	if time.Since(birthdate) < 18*365.25*24*time.Hour {
		return PersonalInformation{}, ErrPersonalInformationUnderage
	}

	return PersonalInformation{
		dniType:   dniType,
		dniNumber: dniNumber,
		firstName: firstName,
		lastName:  lastName,
		birthdate: birthdate,
	}, nil
}

func (v PersonalInformation) DNIType() DNIType     { return v.dniType }
func (v PersonalInformation) DNINumber() string    { return v.dniNumber }
func (v PersonalInformation) FirstName() string    { return v.firstName }
func (v PersonalInformation) LastName() string     { return v.lastName }
func (v PersonalInformation) Birthdate() time.Time { return v.birthdate }

func (v PersonalInformation) Equals(other PersonalInformation) bool {
	return v.dniType == other.dniType &&
		v.dniNumber == other.dniNumber &&
		v.firstName == other.firstName &&
		v.lastName == other.lastName &&
		v.birthdate.Equal(other.birthdate)
}

func (v PersonalInformation) ToDTO() PersonalInformationDTO {
	return PersonalInformationDTO{
		DNIType:   v.dniType,
		DNINumber: v.dniNumber,
		FirstName: v.firstName,
		LastName:  v.lastName,
		Birthdate: v.birthdate.Format(time.RFC3339),
	}
}
