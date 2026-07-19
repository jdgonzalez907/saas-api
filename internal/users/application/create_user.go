package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/users/domain"
)

type CreateUser interface {
	Execute(ctx context.Context, user *domain.User) (*domain.User, error)
}

type createUser struct {
	userRepository domain.UserRepository
}

func NewCreateUser(userRepository domain.UserRepository) CreateUser {
	return &createUser{userRepository: userRepository}
}

func (uc *createUser) Execute(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := uc.userRepository.FindByPhone(ctx, user.Phone())
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if existingUser != nil {
		return nil, uc.wrapError(domain.ErrUserPhoneAlreadyExists)
	}

	if user.Email() != nil {
		existingUser, err = uc.userRepository.FindByEmail(ctx, *user.Email())
		if err != nil {
			return nil, uc.wrapError(err)
		}
		if existingUser != nil {
			return nil, uc.wrapError(domain.ErrUserEmailAlreadyExists)
		}
	}

	if err := uc.userRepository.Create(ctx, user); err != nil {
		return nil, uc.wrapError(err)
	}

	return user, nil
}

func (uc *createUser) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrCreateUser, err)
}
