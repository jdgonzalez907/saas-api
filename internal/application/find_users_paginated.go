package application

import (
	"context"
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type FindUsersPaginatedUseCase interface {
	Execute(ctx context.Context, pagination domain.Pagination) (domain.PaginatedUsers, error)
}

type findUsersPaginatedUseCase struct {
	userRepository domain.UserRepository
}

func NewFindUsersPaginatedUseCase(userRepository domain.UserRepository) FindUsersPaginatedUseCase {
	return &findUsersPaginatedUseCase{
		userRepository: userRepository,
	}
}

func (u *findUsersPaginatedUseCase) Execute(ctx context.Context, pagination domain.Pagination) (domain.PaginatedUsers, error) {
	users, err := u.userRepository.FindAll(ctx, pagination)
	if err != nil {
		return domain.PaginatedUsers{}, fmt.Errorf("%v: %w", domain.ErrFindingUsers, err)
	}

	var nextCursor *int64
	if len(users) > 0 {
		nextCursorVal := users[len(users)-1].ID()
		nextCursor = &nextCursorVal
	}

	return domain.NewPaginatedUsers(users, nextCursor), nil
}
