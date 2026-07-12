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
	userFound, err := u.userRepository.FindById(user.ID())
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPersonalInformation, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	updatedUser, err := userFound.WithPersonalInformation(
		user.Identification(),
		user.FirstName(),
		user.LastName(),
		user.Address(),
		user.BirthDate(),
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
