package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

type RouterParams struct {
	FindUserByID              *FindUserByIDController
	CreateUser                *CreateUserController
	DeleteUser                *DeleteUserController
	UpdatePersonalInformation *UpdateUserPersonalInformationController
	FindUsersPaginated        *FindUsersPaginatedController
	UpdateEmail               *UpdateUserEmailController
	UpdatePhone               *UpdateUserPhoneController
}

func NewRouter(params RouterParams) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(sharedHttp.ErrorLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(sharedHttp.JSONContentTypeMiddleware)

	if params.FindUserByID != nil {
		r.Get("/users/{id}", params.FindUserByID.Handle)
	}
	if params.CreateUser != nil {
		r.Post("/users", params.CreateUser.Handle)
	}
	if params.DeleteUser != nil {
		r.Delete("/users/{id}", params.DeleteUser.Handle)
	}
	if params.UpdatePersonalInformation != nil {
		r.Put("/users/{id}", params.UpdatePersonalInformation.Handle)
	}
	if params.FindUsersPaginated != nil {
		r.Get("/users", params.FindUsersPaginated.Handle)
	}
	if params.UpdateEmail != nil {
		r.Put("/users/{id}/email", params.UpdateEmail.Handle)
	}
	if params.UpdatePhone != nil {
		r.Put("/users/{id}/phone", params.UpdatePhone.Handle)
	}

	return r
}
