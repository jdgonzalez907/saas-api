package controllers

import (
	"net/http"

	"jdgonzalez907/saas-api/internal/users/application"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
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
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.useCase.Execute(r.Context(), id); err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusNoContent, nil)
}
