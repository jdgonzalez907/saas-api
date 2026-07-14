package controllers

import (
	"net/http"

	"jdgonzalez907/saas-api/internal/posts/application"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

type DeletePostController struct {
	useCase application.DeletePostUseCase
}

func NewDeletePostController(useCase application.DeletePostUseCase) *DeletePostController {
	return &DeletePostController{
		useCase: useCase,
	}
}

func (c *DeletePostController) Handle(w http.ResponseWriter, r *http.Request) {
	deletedByID, err := sharedHttp.GetAuthenticatedUserID(r.Context())
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.useCase.Execute(r.Context(), id, deletedByID); err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusNoContent, nil)
}
