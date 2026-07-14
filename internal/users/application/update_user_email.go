package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type UpdateUserEmailUseCase interface {
	Execute(ctx context.Context, id int64, email *domain.Email) error
}

type updateUserEmailUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateUserEmailUseCase(userRepository domain.UserRepository) UpdateUserEmailUseCase {
	return &updateUserEmailUseCase{userRepository: userRepository}
}

func (u *updateUserEmailUseCase) Execute(ctx context.Context, id int64, email *domain.Email) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserEmail, err)
	}

	if userFound == nil {
		return domain.ErrUserNotFound
	}

	currentEmail := userFound.Email()
	emailsEqual := (currentEmail == nil && email == nil) || (currentEmail != nil && email != nil && currentEmail.Equals(*email))
	if emailsEqual {
		return nil
	}

	if email != nil {
		foundEmail, err := u.userRepository.FindByEmail(ctx, *email)
		if err != nil {
			return fmt.Errorf("%v: %w", domain.ErrUpdatingUserEmail, err)
		}
		if foundEmail != nil {
			return domain.ErrUserEmailAlreadyExists
		}
	}

	updatedUser := userFound.WithEmail(email)

	err = u.userRepository.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrUpdatingUserEmail, err)
	}

	return nil
}
