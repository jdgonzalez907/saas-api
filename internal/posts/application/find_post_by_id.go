package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/posts/domain"
)

type FindPostByIDUseCase interface {
	Execute(ctx context.Context, id int64) (*domain.Post, error)
}

type findPostByIDUseCase struct {
	postRepository domain.PostRepository
}

func NewFindPostByIDUseCase(postRepository domain.PostRepository) FindPostByIDUseCase {
	return &findPostByIDUseCase{postRepository: postRepository}
}

func (f *findPostByIDUseCase) Execute(ctx context.Context, id int64) (*domain.Post, error) {
	post, err := f.postRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrFindingPost, err)
	}
	if post == nil {
		return nil, domain.ErrPostNotFound
	}
	return post, nil
}
