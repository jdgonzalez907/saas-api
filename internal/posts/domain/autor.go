package domain

import "errors"

var (
	ErrAutorNotFound         = errors.New("autor not found")
	ErrAutorIDRequired       = errors.New("autor id is required")
	ErrAutorFullNameRequired = errors.New("fullname requerido")
)

type Autor struct {
	id       int64
	fullName string
}

type AutorDTO struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}

func NewAutor(id int64, fullName string) (*Autor, error) {
	if id <= 0 {
		return nil, ErrAutorIDRequired
	}
	if fullName == "" {
		return nil, ErrAutorFullNameRequired
	}
	return &Autor{id: id, fullName: fullName}, nil
}

func (a *Autor) ID() int64        { return a.id }
func (a *Autor) FullName() string { return a.fullName }

func (a *Autor) Equals(other *Autor) bool {
	if other == nil {
		return false
	}
	return a.id == other.id
}

func (a *Autor) ToDTO() AutorDTO {
	return AutorDTO{
		ID:       a.id,
		FullName: a.fullName,
	}
}
