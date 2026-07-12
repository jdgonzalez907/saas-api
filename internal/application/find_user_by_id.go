package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type FindUserByIdUseCase interface {
	Execute(id int) (*domain.User, error)
}

type findUserByIdUseCase struct {
	userRepository domain.UserRepository
}

func NewFindUserByIdUseCase(userRepository domain.UserRepository) FindUserByIdUseCase {
	return &findUserByIdUseCase{
		userRepository: userRepository,
	}
}

func (u *findUserByIdUseCase) Execute(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	userFound, err := u.userRepository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", domain.ErrFindingUserByID, err)
	}

	if userFound == nil {
		return nil, domain.ErrUserNotFound
	}

	return userFound, nil
}
