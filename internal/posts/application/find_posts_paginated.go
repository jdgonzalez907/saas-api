package application

import (
	"context"
	"fmt"
	"time"

	"jdgonzalez907/saas-api/internal/posts/domain"
)

type FindPostsPaginatedUseCase interface {
	Execute(ctx context.Context, status domain.PostStatus, pagination domain.Pagination) (domain.PaginatedPosts, error)
}

type findPostsPaginatedUseCase struct {
	postRepository domain.PostRepository
}

func NewFindPostsPaginatedUseCase(postRepository domain.PostRepository) FindPostsPaginatedUseCase {
	return &findPostsPaginatedUseCase{postRepository: postRepository}
}

func (f *findPostsPaginatedUseCase) Execute(ctx context.Context, status domain.PostStatus, pagination domain.Pagination) (domain.PaginatedPosts, error) {
	posts, err := f.postRepository.FindAll(ctx, status, pagination)
	if err != nil {
		return domain.PaginatedPosts{}, fmt.Errorf("%v: %w", domain.ErrFindingPosts, err)
	}

	var nextPublishedAt *time.Time
	var nextID *int64
	if len(posts) > 0 {
		lastPost := posts[len(posts)-1]
		nextPublishedAt = lastPost.PublishedAt()
		nextIDVal := lastPost.ID()
		nextID = &nextIDVal
	}

	return domain.NewPaginatedPosts(posts, nextPublishedAt, nextID), nil
}
