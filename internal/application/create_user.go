package application

import (
	"fmt"
	"jdgonzalez907/users-api/internal/domain"
)

type CreateUserUseCase interface {
	Execute(user *domain.User) error
}

type createUserUseCase struct {
	userRepository domain.UserRepository
}

func NewCreateUserUseCase(userRepository domain.UserRepository) CreateUserUseCase {
	return &createUserUseCase{userRepository: userRepository}
}

func (c *createUserUseCase) Execute(user *domain.User) error {
	userFound, err := c.userRepository.FindById(user.ID())
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
	}
	if userFound != nil {
		return domain.ErrUserIDAlreadyExists
	}

	userFound, err = c.userRepository.FindByPhone(user.Phone())
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
	}
	if userFound != nil {
		return domain.ErrUserPhoneAlreadyExists
	}

	if user.Email() != nil {
		userFound, err = c.userRepository.FindByEmail(*user.Email())
		if err != nil {
			return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
		}
		if userFound != nil {
			return domain.ErrUserEmailAlreadyExists
		}
	}

	err = c.userRepository.Create(user)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrCreatingUser, err)
	}

	return nil
}
