package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/users/domain"
)

type ChangePhone interface {
	Execute(ctx context.Context, executedBy int64, phone domain.Phone) (*domain.User, error)
}

type changePhone struct {
	userRepository domain.UserRepository
}

func NewChangePhone(userRepository domain.UserRepository) ChangePhone {
	return &changePhone{userRepository: userRepository}
}

func (uc *changePhone) Execute(ctx context.Context, executedBy int64, phone domain.Phone) (*domain.User, error) {
	user, err := uc.userRepository.FindByID(ctx, executedBy)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}

	existingByPhone, err := uc.userRepository.FindByPhone(ctx, phone)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if existingByPhone != nil && !existingByPhone.Equals(user) {
		return nil, uc.wrapError(domain.ErrUserPhoneAlreadyExists)
	}

	if err := user.ChangePhone(phone, executedBy); err != nil {
		return nil, uc.wrapError(err)
	}

	if err := uc.userRepository.Update(ctx, user); err != nil {
		return nil, uc.wrapError(err)
	}

	return user, nil
}

func (uc *changePhone) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrChangePhone, err)
}
