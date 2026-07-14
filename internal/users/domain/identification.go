package domain

import "errors"

var (
	ErrInvalidIdentificationType   = errors.New("identification type must be a standard type")
	ErrInvalidIdentificationNumber = errors.New("identification number cannot be empty")
)

type IdentificationType string

const (
	IDTypeCC       IdentificationType = "CC"
	IDTypeCE       IdentificationType = "CE"
	IDTypePASSPORT IdentificationType = "PASSPORT"
	IDTypeNIT      IdentificationType = "NIT"
	IDTypeIC       IdentificationType = "IC"
)

type Identification struct {
	idType IdentificationType
	number string
}

func NewIdentification(idType IdentificationType, number string) (Identification, error) {
	if number == "" {
		return Identification{}, ErrInvalidIdentificationNumber
	}
	switch idType {
	case IDTypeCC, IDTypeCE, IDTypePASSPORT, IDTypeNIT, IDTypeIC:
		return Identification{
			idType: idType,
			number: number,
		}, nil
	default:
		return Identification{}, ErrInvalidIdentificationType
	}
}

type IdentificationDTO struct {
	Type   IdentificationType `json:"type"`
	Number string             `json:"number"`
}

func (id Identification) Equals(other Identification) bool {
	return id.idType == other.idType && id.number == other.number
}

func (id Identification) ToDTO() IdentificationDTO {
	return IdentificationDTO{
		Type:   id.idType,
		Number: id.number,
	}
}

func (id Identification) Type() IdentificationType {
	return id.idType
}

func (id Identification) Number() string {
	return id.number
}
