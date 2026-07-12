package domain

import "errors"

var ErrInvalidIdentificationType = errors.New("invalid identification type")

type IdentificationType string

const (
	IdType_CC       IdentificationType = "CC"
	IdType_CE       IdentificationType = "CE"
	IdType_PASSPORT IdentificationType = "PASSPORT"
	IdType_NIT      IdentificationType = "NIT"
	IdType_IC       IdentificationType = "IC"
)

type Identification struct {
	idType IdentificationType
	number string
}

func NewIdentification(idType IdentificationType, number string) (Identification, error) {
	switch idType {
	case IdType_CC, IdType_CE, IdType_PASSPORT, IdType_NIT, IdType_IC:
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

func (id Identification) ToDTO() IdentificationDTO {
	return IdentificationDTO{
		Type:   id.idType,
		Number: id.number,
	}
}
