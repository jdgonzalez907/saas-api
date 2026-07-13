package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserEmailController struct {
	useCase application.UpdateUserEmailUseCase
}

func NewUpdateUserEmailController(useCase application.UpdateUserEmailUseCase) *UpdateUserEmailController {
	return &UpdateUserEmailController{
		useCase: useCase,
	}
}

func (c *UpdateUserEmailController) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteIntParam(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body domain.EmailDTO
	if err := ParseJSONBody(r, &body); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var emailPtr *domain.Email
	if body != "" {
		email, err := domain.NewEmail(string(body))
		if err != nil {
			RespondWithDomainError(w, err)
			return
		}
		emailPtr = &email
	}

	err = c.useCase.Execute(id, emailPtr)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
