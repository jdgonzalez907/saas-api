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

type UpdateEmailRequest struct {
	Email *string `json:"email"`
}

func (c *UpdateUserEmailController) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteInt64Param(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body UpdateEmailRequest
	if err := ParseJSONBody(r, &body); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var emailPtr *domain.Email
	if body.Email != nil && *body.Email != "" {
		email, err := domain.NewEmail(*body.Email)
		if err != nil {
			RespondWithDomainError(w, r, err)
			return
		}
		emailPtr = &email
	}

	err = c.useCase.Execute(r.Context(), id, emailPtr)
	if err != nil {
		RespondWithDomainError(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
