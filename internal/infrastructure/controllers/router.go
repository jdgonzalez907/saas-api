package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	findUserByID *FindUserByIDController,
	createUser *CreateUserController,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(JSONContentTypeMiddleware)

	r.Get("/users/{id}", findUserByID.Handle)
	r.Post("/users", createUser.Handle)

	return r
}
