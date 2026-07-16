package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type UpdateUserPersonalInformationUseCase interface {
	Execute(ctx context.Context, id int64, info domain.PersonalInformation) error
}

type updateUserPersonalInformationUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserPersonalInformationUseCase(userRepository domain.UserRepository) UpdateUserPersonalInformationUseCase {
	return &updateUserPersonalInformationUseCase{userRepository: userRepository}
}

func (u *updateUserPersonalInformationUseCase) Execute(ctx context.Context, id int64, info domain.PersonalInformation) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	if info.Equals(userFound.PersonalInformation()) {
		return nil
	}

	updatedUser := userFound.UpdatePersonalInformation(info)

	err = u.userRepository.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	return nil
}
