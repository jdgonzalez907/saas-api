package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
)

type CreatePost interface {
	Execute(ctx context.Context, title, slug, cover string, content []domain.Block, status domain.PostStatus, authorID int64) (*domain.Post, error)
}

type createPost struct {
	postRepository domain.PostRepository
}

func NewCreatePost(postRepository domain.PostRepository) CreatePost {
	return &createPost{postRepository: postRepository}
}

func (uc *createPost) Execute(ctx context.Context, title, slug, cover string, content []domain.Block, status domain.PostStatus, authorID int64) (*domain.Post, error) {
	existingPost, err := uc.postRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if existingPost != nil {
		return nil, uc.wrapError(domain.ErrPostSlugAlreadyExists)
	}

	post, err := domain.New(title, slug, cover, content, status, authorID)
	if err != nil {
		return nil, err
	}

	if err := uc.postRepository.Create(ctx, post); err != nil {
		return nil, uc.wrapError(err)
	}

	return post, nil
}

func (uc *createPost) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrCreatePost, err)
}
