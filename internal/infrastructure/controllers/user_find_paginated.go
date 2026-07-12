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
	limitPtr, err := ParseQueryIntParam(r, "limit")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cursorPtr, err := ParseQueryIntParam(r, "cursor")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	limitVal := 0
	if limitPtr != nil {
		limitVal = *limitPtr
	}

	pagination := domain.NewPagination(cursorPtr, limitVal)

	paginatedUsers, err := c.useCase.Execute(pagination)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, paginatedUsers.ToDTO())
}
