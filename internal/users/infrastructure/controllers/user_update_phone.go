package controllers

import (
	"net/http"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/domain"
)

type UpdateUserPhoneController struct {
	useCase application.UpdateUserPhoneUseCase
}

func NewUpdateUserPhoneController(useCase application.UpdateUserPhoneUseCase) *UpdateUserPhoneController {
	return &UpdateUserPhoneController{
		useCase: useCase,
	}
}

func (c *UpdateUserPhoneController) Handle(w http.ResponseWriter, r *http.Request, _ int64) {
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body domain.PhoneDTO
	if err := sharedHttp.DecodeJSON(r, &body); err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	phone, err := domain.NewPhone(body.CountryCode, body.Number)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	err = c.useCase.Execute(r.Context(), id, phone)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusNoContent, nil)
}
