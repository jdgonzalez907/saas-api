package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterParams struct {
	FindUserByID *FindUserByIDController
	CreateUser   *CreateUserController
	DeleteUser   *DeleteUserController
}

func NewRouter(params RouterParams) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(JSONContentTypeMiddleware)

	if params.FindUserByID != nil {
		r.Get("/users/{id}", params.FindUserByID.Handle)
	}
	if params.CreateUser != nil {
		r.Post("/users", params.CreateUser.Handle)
	}
	if params.DeleteUser != nil {
		r.Delete("/users/{id}", params.DeleteUser.Handle)
	}

	return r
}
