package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
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
	limitPtr, err := ParseQueryInt32Param(r, "limit")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cursorPtr, err := ParseQueryInt64Param(r, "cursor")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	pagination, err := domain.NewPagination(cursorPtr, limitPtr)
	if err != nil {
		RespondWithDomainError(w, r, err)
		return
	}

	paginatedUsers, err := c.useCase.Execute(r.Context(), pagination)
	if err != nil {
		RespondWithDomainError(w, r, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, paginatedUsers.ToDTO())
}
