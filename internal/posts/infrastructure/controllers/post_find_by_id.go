package controllers

import (
	"net/http"

	"jdgonzalez907/saas-api/internal/posts/application"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

type FindPostByIDController struct {
	useCase application.FindPostByIDUseCase
}

func NewFindPostByIDController(useCase application.FindPostByIDUseCase) *FindPostByIDController {
	return &FindPostByIDController{
		useCase: useCase,
	}
}

func (c *FindPostByIDController) Handle(w http.ResponseWriter, r *http.Request, _ int64) {
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusOK, post.ToDTO())
}
