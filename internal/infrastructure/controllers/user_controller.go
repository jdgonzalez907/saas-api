package controllers

import (
	"net/http"

	"jdgonzalez907/users-api/internal/application"
)

type UserController struct {
	findUserByIdUseCase application.FindUserByIdUseCase
}

func NewUserController(findUserByIdUseCase application.FindUserByIdUseCase) *UserController {
	return &UserController{
		findUserByIdUseCase: findUserByIdUseCase,
	}
}

func (h *UserController) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseRouteIntParam(r, "id")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.findUserByIdUseCase.Execute(id)
	if err != nil {
		RespondWithDomainError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, user.ToDTO())
}
