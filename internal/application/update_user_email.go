package application

import (
	"fmt"

	"jdgonzalez907/users-api/internal/domain"
)

type UpdateUserEmailUseCase interface {
	Execute(id int, email *domain.Email) error
}

type updateUserEmailUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserEmailUseCase(userRepository domain.UserRepository) UpdateUserEmailUseCase {
	return &updateUserEmailUseCase{userRepository: userRepository}
}

func (u *updateUserEmailUseCase) Execute(id int, email *domain.Email) error {
	userFound, err := u.userRepository.FindById(id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserEmail, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	if email != nil {
		emailChanged := userFound.Email() == nil || *userFound.Email() != *email
		if emailChanged {
			foundEmail, err := u.userRepository.FindByEmail(*email)
			if err != nil {
				return fmt.Errorf("%v: %w", domain.ErrUpdatingUserEmail, err)
			}
			if foundEmail != nil {
				return domain.ErrUserEmailAlreadyExists
			}
		}
	}

	updatedUser := userFound.WithEmail(email)

	err = u.userRepository.Update(updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserEmail, err)
	}

	return nil
}
