package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
)

type FindPostByID interface {
	Execute(ctx context.Context, id int64) (*domain.Post, error)
}

type findPostByID struct {
	postRepository domain.PostRepository
}

func NewFindPostByID(postRepository domain.PostRepository) FindPostByID {
	return &findPostByID{postRepository: postRepository}
}

func (uc *findPostByID) Execute(ctx context.Context, id int64) (*domain.Post, error) {
	post, err := uc.postRepository.FindByID(ctx, id)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if post == nil {
		return nil, uc.wrapError(domain.ErrPostNotFound)
	}
	return post, nil
}

func (uc *findPostByID) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrFindPostByID, err)
}