package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type FindUserByIDUseCase interface {
	Execute(ctx context.Context, id int64) (*domain.User, error)
}

type findUserByIDUseCase struct {
	userRepository domain.UserRepository
}

func NewFindUserByIDUseCase(userRepository domain.UserRepository) FindUserByIDUseCase {
	return &findUserByIDUseCase{
		userRepository: userRepository,
	}
}

func (u *findUserByIDUseCase) Execute(ctx context.Context, id int64) (*domain.User, error) {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return nil, err
	}

	userFound, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrFindingUserByID, err)
	}

	if userFound == nil {
		return nil, domain.ErrUserNotFound
	}

	return userFound, nil
}
