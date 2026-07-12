package controllers

import (
	"errors"
	"net/http"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
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
		if errors.Is(err, domain.ErrUserNotFound) {
			RespondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidUserID) {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	RespondWithJSON(w, http.StatusOK, user.ToDTO())
}
