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

func (idType IdentificationType) String() string {
	return string(idType)
}

type Identification struct {
	Type   IdentificationType `json:"type"`
	Number string             `json:"number"`
}

func NewIdentification(idType IdentificationType, number string) (Identification, error) {
	switch idType {
	case IdType_CC, IdType_CE, IdType_PASSPORT, IdType_NIT, IdType_IC:
		return Identification{
			Type:   idType,
			Number: number,
		}, nil
	default:
		return Identification{}, ErrInvalidIdentificationType
	}
}

func (id Identification) String() string {
	return id.Type.String() + id.Number
}
