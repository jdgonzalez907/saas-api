package domain

import "context"

type PostRepository interface {
	FindByID(ctx context.Context, id int64) (*Post, error)
	FindAll(ctx context.Context, status PostStatus, pagination Pagination) ([]*Post, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int64, deletedByID int64) error
}
