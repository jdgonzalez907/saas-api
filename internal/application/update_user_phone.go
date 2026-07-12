package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserPhoneUseCase interface {
	Execute(id int, phone domain.Phone) error
}

type updateUserPhoneUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserPhoneUseCase(userRepository domain.UserRepository) UpdateUserPhoneUseCase {
	return &updateUserPhoneUseCase{userRepository: userRepository}
}

func (u *updateUserPhoneUseCase) Execute(id int, phone domain.Phone) error {
	userFound, err := u.userRepository.FindById(id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPhone, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	userFoundDTO := userFound.ToDTO()

	if phone != userFoundDTO.Phone {
		foundPhone, err := u.userRepository.FindByPhone(phone)
		if err != nil {
			return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPhone, err)
		}
		if foundPhone != nil {
			return domain.ErrUserPhoneAlreadyExists
		}
	}

	updatedUser := userFound.WithPhone(phone)

	err = u.userRepository.Update(updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserPhone, err)
	}

	return nil
}
