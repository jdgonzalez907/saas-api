package controllers

import (
	"net/http"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	"jdgonzalez907/saas-api/internal/users/application"
)

type DeleteUserController struct {
	useCase application.DeleteUserUseCase
}

func NewDeleteUserController(useCase application.DeleteUserUseCase) *DeleteUserController {
	return &DeleteUserController{
		useCase: useCase,
	}
}

func (c *DeleteUserController) Handle(w http.ResponseWriter, r *http.Request, userID int64) {
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.useCase.Execute(r.Context(), id, userID); err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusNoContent, nil)
}
