package controllers

import (
	"net/http"

	"jdgonzalez907/saas-api/internal/posts/application"
	"jdgonzalez907/saas-api/internal/posts/domain"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

type CreatePostController struct {
	useCase application.CreatePostUseCase
}

func NewCreatePostController(useCase application.CreatePostUseCase) *CreatePostController {
	return &CreatePostController{
		useCase: useCase,
	}
}

type CreatePostRequest struct {
	domain.ContentInformationDTO
	Status string `json:"status"`
}

func (c *CreatePostController) Handle(w http.ResponseWriter, r *http.Request) {
	authorID, err := sharedHttp.GetAuthenticatedUserID(r.Context())
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	var req CreatePostRequest
	if err := sharedHttp.DecodeJSON(r, &req); err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	contentInfo, err := domain.ContentInformationFromDTO(req.ContentInformationDTO)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	status, err := domain.NewPostStatus(req.Status)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	post, err := domain.NewPostWithoutID(contentInfo, status, authorID)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	if err := c.useCase.Execute(r.Context(), post); err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusCreated, post.ToDTO())
}
