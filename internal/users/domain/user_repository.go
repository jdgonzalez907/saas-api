package domain

import "context"

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByPhone(ctx context.Context, phone Phone) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	FindAll(ctx context.Context, pagination Pagination) ([]*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
