package application

import (
	"context"
	"fmt"
	"jdgonzalez907/saas-api/internal/users/domain"
)

type DeleteUserUseCase interface {
	Execute(ctx context.Context, id int64) error
}

type deleteUserUseCase struct {
	userRepository domain.UserRepository
}

func NewDeleteUserUseCase(userRepository domain.UserRepository) DeleteUserUseCase {
	return &deleteUserUseCase{userRepository: userRepository}
}

func (d *deleteUserUseCase) Execute(ctx context.Context, id int64) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := d.userRepository.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrDeletingUser, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	err = d.userRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrDeletingUser, err)
	}

	return nil
}
