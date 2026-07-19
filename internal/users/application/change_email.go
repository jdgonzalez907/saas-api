package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/users/domain"
)

type ChangeEmail interface {
	Execute(ctx context.Context, executedBy int64, email domain.Email) (*domain.User, error)
}

type changeEmail struct {
	userRepository domain.UserRepository
}

func NewChangeEmail(userRepository domain.UserRepository) ChangeEmail {
	return &changeEmail{userRepository: userRepository}
}

func (uc *changeEmail) Execute(ctx context.Context, executedBy int64, email domain.Email) (*domain.User, error) {
	user, err := uc.userRepository.FindByID(ctx, executedBy)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}

	existingByEmail, err := uc.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if existingByEmail != nil && !existingByEmail.Equals(user) {
		return nil, uc.wrapError(domain.ErrUserEmailAlreadyExists)
	}

	if err := user.ChangeEmail(email, executedBy); err != nil {
		return nil, uc.wrapError(err)
	}

	if err := uc.userRepository.Update(ctx, user); err != nil {
		return nil, uc.wrapError(err)
	}

	return user, nil
}

func (uc *changeEmail) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrChangeEmail, err)
}
