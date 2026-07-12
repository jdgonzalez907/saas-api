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
	userFound, err := u.userRepository.FindById(user.ID)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	if user.Phone != userFound.Phone {
		foundPhone, err := u.userRepository.FindByPhone(user.Phone)
		if err != nil {
			return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
		}
		if foundPhone != nil {
			return domain.ErrUserPhoneAlreadyExists
		}
	}

	if user.Email != nil {
		emailChanged := userFound.Email == nil || *userFound.Email != *user.Email
		if emailChanged {
			foundEmail, err := u.userRepository.FindByEmail(*user.Email)
			if err != nil {
				return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
			}
			if foundEmail != nil {
				return domain.ErrUserEmailAlreadyExists
			}
		}
	}

	err = u.userRepository.Update(user)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUser, err)
	}

	return nil
}
