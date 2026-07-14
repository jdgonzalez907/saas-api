package controllers

import (
	"net/http"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/domain"
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
	if err := sharedHttp.DecodeJSON(r, &body); err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	identification, err := domain.NewIdentification(body.Identification.Type, body.Identification.Number)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	phone, err := domain.NewPhone(body.Phone.CountryCode, body.Phone.Number)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	var email *domain.Email
	if body.Email != nil {
		e, err := domain.NewEmail(string(*body.Email))
		if err != nil {
			sharedHttp.RespondWithDomainError(w, r, err)
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
			sharedHttp.RespondWithDomainError(w, r, err)
			return
		}
		address = &a
	}

	var birthDate *domain.BirthDate
	if body.BirthDate != nil {
		bd, err := domain.NewBirthDate(string(*body.BirthDate))
		if err != nil {
			sharedHttp.RespondWithDomainError(w, r, err)
			return
		}
		birthDate = &bd
	}

	personalInfo, err := domain.NewPersonalInformation(identification, body.FirstName, body.LastName, address, birthDate)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	user, _ := domain.NewUserWithoutID(personalInfo, phone, email)

	if err := c.useCase.Execute(r.Context(), user); err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusCreated, user.ToDTO())
}
