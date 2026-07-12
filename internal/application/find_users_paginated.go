package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type FindUsersPaginatedUseCase interface {
	Execute(pagination domain.Pagination) (domain.PaginatedUsersDTO, error)
}

type findUsersPaginatedUseCase struct {
	userRepository domain.UserRepository
}

func NewFindUsersPaginatedUseCase(userRepository domain.UserRepository) FindUsersPaginatedUseCase {
	return &findUsersPaginatedUseCase{
		userRepository: userRepository,
	}
}

func (u *findUsersPaginatedUseCase) Execute(pagination domain.Pagination) (domain.PaginatedUsersDTO, error) {

	users, err := u.userRepository.FindAll(pagination)
	if err != nil {
		return domain.PaginatedUsersDTO{}, fmt.Errorf("%v: %w", domain.ErrFindingUsers, err)
	}

	userDTOs := make([]domain.UserDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, *user.ToDTO())
	}

	var nextCursor *int
	if len(users) > 0 {
		nextCursorVal := users[len(users)-1].ToDTO().ID
		nextCursor = &nextCursorVal
	}

	return domain.PaginatedUsersDTO{
		Users:      userDTOs,
		NextCursor: nextCursor,
	}, nil
}
