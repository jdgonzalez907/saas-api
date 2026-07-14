package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
)

type FindUserByIDController struct {
	useCase application.FindUserByIDUseCase
}

func NewFindUserByIDController(useCase application.FindUserByIDUseCase) *FindUserByIDController {
	return &FindUserByIDController{
		useCase: useCase,
	}
}

func (c *FindUserByIDController) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteInt64Param(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		RespondWithDomainError(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, user.ToDTO())
}
