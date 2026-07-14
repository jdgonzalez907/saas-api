package controllers

import (
	"net/http"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/domain"
)

type FindUsersPaginatedController struct {
	useCase application.FindUsersPaginatedUseCase
}

func NewFindUsersPaginatedController(useCase application.FindUsersPaginatedUseCase) *FindUsersPaginatedController {
	return &FindUsersPaginatedController{
		useCase: useCase,
	}
}

func (c *FindUsersPaginatedController) Handle(w http.ResponseWriter, r *http.Request) {
	limitPtr, err := sharedHttp.ParseQueryInt32Param(r, "limit")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cursorPtr, err := sharedHttp.ParseQueryInt64Param(r, "cursor")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	pagination, err := domain.NewPagination(cursorPtr, limitPtr)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	paginatedUsers, err := c.useCase.Execute(r.Context(), pagination)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusOK, paginatedUsers.ToDTO())
}
