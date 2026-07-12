package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserUseCase interface {
	Execute(user *domain.User) error
}

type updateUserUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserUseCase(userRepository domain.UserRepository) UpdateUserUseCase {
	return &updateUserUseCase{userRepository: userRepository}
}

func (u *updateUserUseCase) Execute(user *domain.User) error {
	userDTO := user.ToDTO()

	userFound, err := u.userRepository.FindById(userDTO.ID)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	userFoundDTO := userFound.ToDTO()

	if userDTO.Phone != userFoundDTO.Phone {
		foundPhone, err := u.userRepository.FindByPhone(userDTO.Phone)
		if err != nil {
			return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
		}
		if foundPhone != nil {
			return domain.ErrUserPhoneAlreadyExists
		}
	}

	if userDTO.Email != nil {
		emailChanged := userFoundDTO.Email == nil || *userFoundDTO.Email != *userDTO.Email
		if emailChanged {
			foundEmail, err := u.userRepository.FindByEmail(*userDTO.Email)
			if err != nil {
				return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
			}
			if foundEmail != nil {
				return domain.ErrUserEmailAlreadyExists
			}
		}
	}

	updatedUser, err := userFound.With(domain.UserParams{
		Identification: userDTO.Identification,
		FirstName:      userDTO.FirstName,
		LastName:       userDTO.LastName,
		Phone:          userDTO.Phone,
		Email:          userDTO.Email,
		Address:        userDTO.Address,
		BirthDate:      userDTO.BirthDate,
	})
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
	}

	err = u.userRepository.Update(updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
	}

	return nil
}
