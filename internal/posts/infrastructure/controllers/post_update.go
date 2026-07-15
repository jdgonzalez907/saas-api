package controllers

import (
	"net/http"

	"jdgonzalez907/saas-api/internal/posts/application"
	"jdgonzalez907/saas-api/internal/posts/domain"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
)

type UpdatePostController struct {
	useCase application.UpdatePostUseCase
}

func NewUpdatePostController(useCase application.UpdatePostUseCase) *UpdatePostController {
	return &UpdatePostController{
		useCase: useCase,
	}
}

type UpdatePostRequest struct {
	domain.ContentInformationDTO
	Status string `json:"status"`
}

func (c *UpdatePostController) Handle(w http.ResponseWriter, r *http.Request, userID int64) {
	id, err := sharedHttp.ParseRouteInt64Param(r, "id")
	if err != nil {
		sharedHttp.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req UpdatePostRequest
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

	post, err := c.useCase.Execute(r.Context(), id, contentInfo, status, userID)
	if err != nil {
		sharedHttp.RespondWithDomainError(w, r, err)
		return
	}

	sharedHttp.RespondWithJSON(w, http.StatusOK, post.ToDTO())
}
