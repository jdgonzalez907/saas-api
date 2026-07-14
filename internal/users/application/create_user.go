package application

import (
	"context"
	"fmt"
	"jdgonzalez907/saas-api/internal/users/domain"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, user *domain.User) error
}

type createUserUseCase struct {
	userRepository domain.UserRepository
}

func NewCreateUserUseCase(userRepository domain.UserRepository) CreateUserUseCase {
	return &createUserUseCase{userRepository: userRepository}
}

func (c *createUserUseCase) Execute(ctx context.Context, user *domain.User) error {
	userFound, err := c.userRepository.FindById(ctx, user.ID())
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
	}
	if userFound != nil {
		return domain.ErrUserIDAlreadyExists
	}

	userFound, err = c.userRepository.FindByPhone(ctx, user.Phone())
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
	}
	if userFound != nil {
		return domain.ErrUserPhoneAlreadyExists
	}

	if user.Email() != nil {
		userFound, err = c.userRepository.FindByEmail(ctx, *user.Email())
		if err != nil {
			return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
		}
		if userFound != nil {
			return domain.ErrUserEmailAlreadyExists
		}
	}

	err = c.userRepository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
	}

	return nil
}
