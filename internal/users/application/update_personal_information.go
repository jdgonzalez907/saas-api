package application

import (
	"context"
	"fmt"

	"github.com/jdgonzalez907/saas-api/internal/users/domain"
)

type UpdatePersonalInformation interface {
	Execute(ctx context.Context, executedBy int64, personalInformation domain.PersonalInformation) (*domain.User, error)
}

type updatePersonalInformation struct {
	userRepository domain.UserRepository
}

func NewUpdatePersonalInformation(userRepository domain.UserRepository) UpdatePersonalInformation {
	return &updatePersonalInformation{userRepository: userRepository}
}

func (uc *updatePersonalInformation) Execute(ctx context.Context, executedBy int64, personalInformation domain.PersonalInformation) (*domain.User, error) {
	user, err := uc.userRepository.FindByID(ctx, executedBy)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}

	if err = user.UpdatePersonalInformation(personalInformation, executedBy); err != nil {
		return nil, uc.wrapError(err)
	}

	if err = uc.userRepository.Update(ctx, user); err != nil {
		return nil, uc.wrapError(err)
	}

	return user, nil
}

func (uc *updatePersonalInformation) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrUpdatePersonalInformation, err)
}
