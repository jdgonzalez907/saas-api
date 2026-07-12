package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
)

type DeleteUserController struct {
	useCase application.DeleteUserUseCase
}

func NewDeleteUserController(useCase application.DeleteUserUseCase) *DeleteUserController {
	return &DeleteUserController{
		useCase: useCase,
	}
}

func (c *DeleteUserController) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteIntParam(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.useCase.Execute(id); err != nil {
		RespondWithDomainError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
