package controllers

import (
	"net/http"

	"jdgonzalez907/saas-api/internal/posts/application"
	"jdgonzalez907/saas-api/internal/posts/domain"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
)

type FindPostsPaginatedController struct {
	useCase application.FindPostsPaginatedUseCase
}

func NewFindPostsPaginatedController(useCase application.FindPostsPaginatedUseCase) *FindPostsPaginatedController {
	return &FindPostsPaginatedController{
		useCase: useCase,
	}
}

func (c *FindPostsPaginatedController) Handle(w http.ResponseWriter, r *http.Request, _ int64) {
	statusStr := r.URL.Query().Get("status")
	status, err := domain.NewPostStatus(statusStr)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	limitPtr, err := sharedHttp.ParseQueryInt32Param(r, "limit")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	lastIDPtr, err := sharedHttp.ParseQueryInt64Param(r, "lastID")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	lastPublishedAtPtr, err := sharedHttp.ParseQueryTimeParam(r, "lastPublishedAt")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	pagination, err := domain.NewPagination(lastPublishedAtPtr, lastIDPtr, limitPtr)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	paginatedPosts, err := c.useCase.Execute(r.Context(), status, pagination)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusOK, paginatedPosts.ToDTO())
}
