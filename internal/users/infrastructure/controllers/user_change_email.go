package controllers

import (
	"net/http"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/domain"
)

type ChangeUserEmailController struct {
	useCase application.ChangeUserEmailUseCase
}

func NewChangeUserEmailController(useCase application.ChangeUserEmailUseCase) *ChangeUserEmailController {
	return &ChangeUserEmailController{
		useCase: useCase,
	}
}

type ChangeEmailRequest struct {
	Email *string `json:"email"`
}

func (c *ChangeUserEmailController) Handle(w http.ResponseWriter, r *http.Request, userID int64) {
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var body ChangeEmailRequest
	if err := sharedHttp.DecodeJSON(r, &body); err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var emailPtr *domain.Email
	if body.Email != nil && *body.Email != "" {
		email, err := domain.NewEmail(*body.Email)
		if err != nil {
			sharedHttp.RespondWithDomainError(w, r, err)
			return
		}
		emailPtr = &email
	}

	err = c.useCase.Execute(r.Context(), id, emailPtr, userID)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusNoContent, nil)
}
