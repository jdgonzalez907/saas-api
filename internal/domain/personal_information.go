package domain

type PersonalInformation struct {
	identification Identification
	firstName      string
	lastName       string
	address        *Address
	birthDate      *BirthDate
}

func NewPersonalInformation(
	identification Identification,
	firstName, lastName string,
	address *Address,
	birthDate *BirthDate,
) (PersonalInformation, error) {
	if firstName == "" {
		return PersonalInformation{}, ErrInvalidFirstName
	}
	if lastName == "" {
		return PersonalInformation{}, ErrInvalidLastName
	}

	return PersonalInformation{
		identification: identification,
		firstName:      firstName,
		lastName:       lastName,
		address:        address,
		birthDate:      birthDate,
	}, nil
}

type PersonalInformationDTO struct {
	Identification IdentificationDTO `json:"identification"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	Address        *AddressDTO       `json:"address"`
	BirthDate      *BirthDateDTO     `json:"birth_date"`
}

func (p PersonalInformation) ToDTO() PersonalInformationDTO {
	var addressDTO *AddressDTO
	if p.address != nil {
		dto := p.address.ToDTO()
		addressDTO = &dto
	}

	var birthDateDTO *BirthDateDTO
	if p.birthDate != nil {
		dto := p.birthDate.ToDTO()
		birthDateDTO = &dto
	}

	return PersonalInformationDTO{
		Identification: p.identification.ToDTO(),
		FirstName:      p.firstName,
		LastName:       p.lastName,
		Address:        addressDTO,
		BirthDate:      birthDateDTO,
	}
}
