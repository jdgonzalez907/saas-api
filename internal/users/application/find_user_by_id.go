package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/users/domain"
)

type FindUserByID interface {
	Execute(ctx context.Context, id int64) (*domain.User, error)
}

type findUserID struct {
	userRepository domain.UserRepository
}

func NewFindUserByID(userRepository domain.UserRepository) FindUserByID {
	return &findUserID{userRepository: userRepository}
}

func (uc *findUserID) Execute(ctx context.Context, id int64) (*domain.User, error) {
	user, err := uc.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}
	return user, nil
}

func (uc *findUserID) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrFindUserByID, err)
}