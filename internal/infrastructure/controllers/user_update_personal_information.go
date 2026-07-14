package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserPersonalInformationController struct {
	useCase application.UpdateUserPersonalInformationUseCase
}

func NewUpdateUserPersonalInformationController(
	useCase application.UpdateUserPersonalInformationUseCase,
) *UpdateUserPersonalInformationController {
	return &UpdateUserPersonalInformationController{
		useCase: useCase,
	}
}

func (c *UpdateUserPersonalInformationController) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteIntParam(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body domain.PersonalInformationDTO
	if err := ParseJSONBody(r, &body); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	identification, err := domain.NewIdentification(body.Identification.Type, body.Identification.Number)
	if err != nil {
		RespondWithDomainError(w, r, err)
		return
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
			RespondWithDomainError(w, r, err)
			return
		}
		address = &a
	}

	var birthDate *domain.BirthDate
	if body.BirthDate != nil {
		bd, err := domain.NewBirthDate(string(*body.BirthDate))
		if err != nil {
			RespondWithDomainError(w, r, err)
			return
		}
		birthDate = &bd
	}

	personalInfo, err := domain.NewPersonalInformation(
		identification,
		body.FirstName,
		body.LastName,
		address,
		birthDate,
	)
	if err != nil {
		RespondWithDomainError(w, r, err)
		return
	}

	if err := c.useCase.Execute(r.Context(), id, personalInfo); err != nil {
		RespondWithDomainError(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
