package domain

import "context"

type PostRepository interface {
	FindByID(ctx context.Context, id int64) (*Post, error)
	FindBySlug(ctx context.Context, slug string) (*Post, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
}
