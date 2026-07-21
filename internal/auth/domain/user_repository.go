package domain

import "context"

type UserRepository interface {
	FindByPhone(ctx context.Context, phoneNumber string) (*User, error)
}
