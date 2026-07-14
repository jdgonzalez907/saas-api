package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/posts/domain"
)

type DeletePostUseCase interface {
	Execute(ctx context.Context, id int64, deletedByID int64) error
}

type deletePostUseCase struct {
	postRepository domain.PostRepository
}

func NewDeletePostUseCase(postRepository domain.PostRepository) DeletePostUseCase {
	return &deletePostUseCase{postRepository: postRepository}
}

func (d *deletePostUseCase) Execute(ctx context.Context, id int64, deletedByID int64) error {
	post, err := d.postRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrDeletingPost, err)
	}
	if post == nil {
		return domain.ErrPostNotFound
	}

	err = d.postRepository.Delete(ctx, id, deletedByID)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrDeletingPost, err)
	}

	return nil
}
