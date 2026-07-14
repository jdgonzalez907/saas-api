package controllers

import (
	"net/http"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
	"jdgonzalez907/saas-api/internal/users/application"
)

type FindUserByIDController struct {
	useCase application.FindUserByIDUseCase
}

func NewFindUserByIDController(useCase application.FindUserByIDUseCase) *FindUserByIDController {
	return &FindUserByIDController{
		useCase: useCase,
	}
}

func (c *FindUserByIDController) Handle(w http.ResponseWriter, r *http.Request, _ int64) {
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusOK, user.ToDTO())
}
