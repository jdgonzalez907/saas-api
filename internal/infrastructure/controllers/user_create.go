package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
)

type CreateUserController struct {
	useCase application.CreateUserUseCase
}

func NewCreateUserController(useCase application.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{
		useCase: useCase,
	}
}

func (c *CreateUserController) Handle(w http.ResponseWriter, r *http.Request) {
	var body domain.UserDTO
	if err := ParseJSONBody(r, &body); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	identification, err := domain.NewIdentification(body.Identification.Type, body.Identification.Number)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	phone, err := domain.NewPhone(body.Phone.CountryCode, body.Phone.Number)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	var email *domain.Email
	if body.Email != nil {
		e, err := domain.NewEmail(body.Email.Value)
		if err != nil {
			RespondWithDomainError(w, err)
			return
		}
		email = &e
	}

	var address *domain.Address
	if body.Address != nil {
		a, err := domain.NewAddress(
			body.Address.Street,
			body.Address.City,
			body.Address.State,
			body.Address.Country,
			body.Address.PostalCode,
			body.Address.Description,
		)
		if err != nil {
			RespondWithDomainError(w, err)
			return
		}
		address = &a
	}

	var birthDate *domain.BirthDate
	if body.BirthDate != nil {
		bd, err := domain.NewBirthDate(body.BirthDate.Value)
		if err != nil {
			RespondWithDomainError(w, err)
			return
		}
		birthDate = &bd
	}

	user, err := domain.NewUserWithoutId(identification, body.FirstName, body.LastName, phone, email, address, birthDate)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	if err := c.useCase.Execute(user); err != nil {
		RespondWithDomainError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, nil)
}
