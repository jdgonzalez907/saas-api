package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/posts/domain"
)

type UpdatePostUseCase interface {
	Execute(ctx context.Context, id int64, contentInfo domain.ContentInformation, status domain.PostStatus, lastEditorID int64) (*domain.Post, error)
}

type updatePostUseCase struct {
	postRepository domain.PostRepository
}

func NewUpdatePostUseCase(postRepository domain.PostRepository) UpdatePostUseCase {
	return &updatePostUseCase{postRepository: postRepository}
}

func (u *updatePostUseCase) Execute(ctx context.Context, id int64, contentInfo domain.ContentInformation, status domain.PostStatus, lastEditorID int64) (*domain.Post, error) {
	post, err := u.postRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrUpdatingPost, err)
	}
	if post == nil {
		return nil, domain.ErrPostNotFound
	}

	updatedPost, err := post.WithContentAndStatus(contentInfo, status, lastEditorID)
	if err != nil {
		return nil, err
	}

	err = u.postRepository.Update(ctx, updatedPost)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrUpdatingPost, err)
	}

	return updatedPost, nil
}
