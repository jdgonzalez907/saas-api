package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type UpdateUserPersonalInformationUseCase interface {
	Execute(ctx context.Context, id int64, info domain.PersonalInformation, userID int64) error
}

type updateUserPersonalInformationUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserPersonalInformationUseCase(userRepository domain.UserRepository) UpdateUserPersonalInformationUseCase {
	return &updateUserPersonalInformationUseCase{userRepository: userRepository}
}

func (u *updateUserPersonalInformationUseCase) Execute(ctx context.Context, id int64, info domain.PersonalInformation, userID int64) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingPersonalInformation, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	if info.Equals(userFound.PersonalInformation()) {
		return nil
	}

	updatedUser, err := userFound.UpdatePersonalInformation(info, userID)
	if err != nil {
		return err
	}

	err = u.userRepository.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingPersonalInformation, err)
	}

	return nil
}
