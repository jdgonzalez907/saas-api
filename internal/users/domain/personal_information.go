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

	dniNumberMinLength = 8
	dniNumberMaxLength = 20

	nameMinLength = 4
	nameMaxLength = 100

	minimumAge = 18
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
	if err := validatePersonalInformationDNIType(dniType); err != nil {
		return PersonalInformation{}, err
	}
	if err := validatePersonalInformationDNINumber(dniNumber); err != nil {
		return PersonalInformation{}, err
	}
	if err := validatePersonalInformationFirstName(firstName); err != nil {
		return PersonalInformation{}, err
	}
	if err := validatePersonalInformationLastName(lastName); err != nil {
		return PersonalInformation{}, err
	}
	if err := validatePersonalInformationBirthdate(birthdate); err != nil {
		return PersonalInformation{}, err
	}
	if err := validatePersonalInformationAge(birthdate); err != nil {
		return PersonalInformation{}, err
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

func validatePersonalInformationDNIType(dniType DNIType) error {
	switch dniType {
	case DNITypeCC, DNITypeCE, DNITypeNIT, DNITypePP:
		return nil
	default:
		return ErrPersonalInformationInvalidDNIType
	}
}

func validatePersonalInformationDNINumber(dniNumber string) error {
	if len(dniNumber) < dniNumberMinLength || len(dniNumber) > dniNumberMaxLength {
		return ErrPersonalInformationInvalidDNINumber
	}
	return nil
}

func validatePersonalInformationFirstName(firstName string) error {
	if len(firstName) < nameMinLength || len(firstName) > nameMaxLength {
		return ErrPersonalInformationInvalidFirstName
	}
	return nil
}

func validatePersonalInformationLastName(lastName string) error {
	if len(lastName) < nameMinLength || len(lastName) > nameMaxLength {
		return ErrPersonalInformationInvalidLastName
	}
	return nil
}

func validatePersonalInformationBirthdate(birthdate time.Time) error {
	if birthdate.IsZero() {
		return ErrPersonalInformationEmptyBirthdate
	}
	return nil
}

func validatePersonalInformationAge(birthdate time.Time) error {
	age := time.Since(birthdate)
	oneYear := 365 * 24 * time.Hour
	if age < time.Duration(minimumAge)*oneYear {
		return ErrPersonalInformationUnderage
	}
	return nil
}
