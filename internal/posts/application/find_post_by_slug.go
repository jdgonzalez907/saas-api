package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
)

type FindPostBySlug interface {
	Execute(ctx context.Context, slug string) (*domain.Post, error)
}

type findPostBySlug struct {
	postRepository domain.PostRepository
}

func NewFindPostBySlug(postRepository domain.PostRepository) FindPostBySlug {
	return &findPostBySlug{postRepository: postRepository}
}

func (uc *findPostBySlug) Execute(ctx context.Context, slug string) (*domain.Post, error) {
	post, err := uc.postRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if post == nil {
		return nil, uc.wrapError(domain.ErrPostNotFound)
	}
	return post, nil
}

func (uc *findPostBySlug) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrFindPostBySlug, err)
}