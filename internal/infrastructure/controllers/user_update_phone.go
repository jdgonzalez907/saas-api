package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserPhoneController struct {
	useCase application.UpdateUserPhoneUseCase
}

func NewUpdateUserPhoneController(useCase application.UpdateUserPhoneUseCase) *UpdateUserPhoneController {
	return &UpdateUserPhoneController{
		useCase: useCase,
	}
}

func (c *UpdateUserPhoneController) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteIntParam(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body domain.PhoneDTO
	if err := ParseJSONBody(r, &body); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	phone, err := domain.NewPhone(body.CountryCode, body.Number)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	err = c.useCase.Execute(r.Context(), id, phone)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
