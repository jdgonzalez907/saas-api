package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/posts/domain"
)

type CreatePostUseCase interface {
	Execute(ctx context.Context, post *domain.Post) error
}

type createPostUseCase struct {
	postRepository domain.PostRepository
}

func NewCreatePostUseCase(postRepository domain.PostRepository) CreatePostUseCase {
	return &createPostUseCase{postRepository: postRepository}
}
func (c *createPostUseCase) Execute(ctx context.Context, post *domain.Post) error {
	err := c.postRepository.Create(ctx, post)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingPost, err)
	}

	return nil
}
