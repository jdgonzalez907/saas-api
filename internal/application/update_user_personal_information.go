package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserPersonalInformationUseCase interface {
	Execute(id int, info domain.PersonalInformation) error
}

type updateUserPersonalInformationUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserPersonalInformationUseCase(userRepository domain.UserRepository) UpdateUserPersonalInformationUseCase {
	return &updateUserPersonalInformationUseCase{userRepository: userRepository}
}

func (u *updateUserPersonalInformationUseCase) Execute(id int, info domain.PersonalInformation) error {
	userFound, err := u.userRepository.FindById(id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	updatedUser := userFound.WithPersonalInformation(info)

	err = u.userRepository.Update(updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	return nil
}
