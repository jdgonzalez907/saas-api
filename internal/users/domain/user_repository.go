package domain

import "context"

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	FindByPhone(ctx context.Context, phone Phone) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
