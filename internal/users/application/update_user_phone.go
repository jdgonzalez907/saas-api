package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type UpdateUserPhoneUseCase interface {
	Execute(ctx context.Context, id int64, phone domain.Phone) error
}

type updateUserPhoneUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserPhoneUseCase(userRepository domain.UserRepository) UpdateUserPhoneUseCase {
	return &updateUserPhoneUseCase{userRepository: userRepository}
}

func (u *updateUserPhoneUseCase) Execute(ctx context.Context, id int64, phone domain.Phone) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := u.userRepository.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPhone, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	if phone.Equals(userFound.Phone()) {
		return nil
	}

	foundPhone, err := u.userRepository.FindByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPhone, err)
	}
	if foundPhone != nil {
		return domain.ErrUserPhoneAlreadyExists
	}

	updatedUser := userFound.WithPhone(phone)

	err = u.userRepository.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPhone, err)
	}

	return nil
}
