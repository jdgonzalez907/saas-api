package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserPersonalInformationUseCase interface {
	Execute(user *domain.User) error
}

type updateUserPersonalInformationUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserPersonalInformationUseCase(userRepository domain.UserRepository) UpdateUserPersonalInformationUseCase {
	return &updateUserPersonalInformationUseCase{userRepository: userRepository}
}

func (u *updateUserPersonalInformationUseCase) Execute(user *domain.User) error {
	userDTO := user.ToDTO()

	userFound, err := u.userRepository.FindById(userDTO.ID)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	updatedUser, err := userFound.WithPersonalInformation(
		userDTO.Identification,
		userDTO.FirstName,
		userDTO.LastName,
		userDTO.Address,
		userDTO.BirthDate,
	)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	err = u.userRepository.Update(updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	return nil
}
