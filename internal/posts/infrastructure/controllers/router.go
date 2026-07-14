package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

type RouterParams struct {
	CreatePost         *CreatePostController
	FindPostByID       *FindPostByIDController
	UpdatePost         *UpdatePostController
	DeletePost         *DeletePostController
	FindPostsPaginated *FindPostsPaginatedController
}

func NewRouter(params RouterParams) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(sharedHttp.ErrorLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(sharedHttp.JSONContentTypeMiddleware)

	if params.CreatePost != nil {
		r.Post("/posts", sharedHttp.Protected(params.CreatePost.Handle))
	}
	if params.FindPostByID != nil {
		r.Get("/posts/{id}", sharedHttp.Protected(params.FindPostByID.Handle))
	}
	if params.UpdatePost != nil {
		r.Put("/posts/{id}", sharedHttp.Protected(params.UpdatePost.Handle))
	}
	if params.DeletePost != nil {
		r.Delete("/posts/{id}", sharedHttp.Protected(params.DeletePost.Handle))
	}
	if params.FindPostsPaginated != nil {
		r.Get("/posts", sharedHttp.Protected(params.FindPostsPaginated.Handle))
	}

	return r
}
