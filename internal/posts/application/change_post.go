package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/posts/domain"
)

type ChangePostUseCase interface {
	Execute(ctx context.Context, id int64, contentInfo domain.ContentInformation, status domain.PostStatus, authorID int64) (*domain.Post, error)
}

type changePostUseCase struct {
	postRepository domain.PostRepository
}

func NewChangePostUseCase(postRepository domain.PostRepository) ChangePostUseCase {
	return &changePostUseCase{postRepository: postRepository}
}

func (u *changePostUseCase) Execute(ctx context.Context, id int64, contentInfo domain.ContentInformation, status domain.PostStatus, authorID int64) (*domain.Post, error) {
	post, err := u.postRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrChangingPost, err)
	}
	if post == nil {
		return nil, domain.ErrPostNotFound
	}

	updatedPost, err := post.UpdateContentAndStatus(contentInfo, status, authorID)
	if err != nil {
		return nil, err
	}

	err = u.postRepository.Update(ctx, updatedPost)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrChangingPost, err)
	}

	return updatedPost, nil
}
