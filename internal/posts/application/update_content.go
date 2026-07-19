package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
)

type UpdateContent interface {
	Execute(ctx context.Context, id, executedBy int64, title, slug, cover string, content []domain.Block, status domain.PostStatus) (*domain.Post, error)
}

type updateContent struct {
	postRepository domain.PostRepository
}

func NewUpdateContent(postRepository domain.PostRepository) UpdateContent {
	return &updateContent{postRepository: postRepository}
}

func (uc *updateContent) Execute(ctx context.Context, id, executedBy int64, title, slug, cover string, content []domain.Block, status domain.PostStatus) (*domain.Post, error) {
	post, err := uc.postRepository.FindByID(ctx, id)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if post == nil {
		return nil, uc.wrapError(domain.ErrPostNotFound)
	}

	if post.Slug() != slug {
		existingBySlug, err := uc.postRepository.FindBySlug(ctx, slug)
		if err != nil {
			return nil, uc.wrapError(err)
		}
		if existingBySlug != nil {
			return nil, uc.wrapError(domain.ErrPostSlugAlreadyExists)
		}
	}

	if err := post.UpdateContent(title, slug, cover, content, status, executedBy); err != nil {
		return nil, err
	}

	if err := uc.postRepository.Update(ctx, post); err != nil {
		return nil, uc.wrapError(err)
	}

	return post, nil
}

func (uc *updateContent) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrUpdateContent, err)
}
