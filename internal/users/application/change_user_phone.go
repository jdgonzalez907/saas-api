package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type ChangeUserPhoneUseCase interface {
	Execute(ctx context.Context, id int64, phone domain.Phone) error
}

type changeUserPhoneUseCase struct {
	userRepository domain.UserRepository
}

func NewChangeUserPhoneUseCase(userRepository domain.UserRepository) ChangeUserPhoneUseCase {
	return &changeUserPhoneUseCase{userRepository: userRepository}
}

func (u *changeUserPhoneUseCase) Execute(ctx context.Context, id int64, phone domain.Phone) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingPhone, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	if phone.Equals(userFound.Phone()) {
		return nil
	}

	foundPhone, err := u.userRepository.FindByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingPhone, err)
	}
	if foundPhone != nil {
		return domain.ErrUserPhoneAlreadyExists
	}

	updatedUser := userFound.ChangePhone(phone)

	err = u.userRepository.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingPhone, err)
	}

	return nil
}
