package application

import (
	"fmt"
	"jdgonzalez907/users-api/internal/domain"
)

type DeleteUserUseCase interface {
	Execute(id string) error
}

type deleteUserUseCase struct {
	userRepository domain.UserRepository
}

func NewDeleteUserUseCase(userRepository domain.UserRepository) DeleteUserUseCase {
	return &deleteUserUseCase{userRepository: userRepository}
}

func (d *deleteUserUseCase) Execute(id string) error {
	userFound, err := d.userRepository.FindById(id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrDeletingUser, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	err = d.userRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrDeletingUser, err)
	}

	return nil
}
