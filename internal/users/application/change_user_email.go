package application

import (
	"context"
	"fmt"

	"jdgonzalez907/saas-api/internal/users/domain"
)

type ChangeUserEmailUseCase interface {
	Execute(ctx context.Context, id int64, email *domain.Email, userID int64) error
}

type changeUserEmailUseCase struct {
	userRepository domain.UserRepository
}

func NewChangeUserEmailUseCase(userRepository domain.UserRepository) ChangeUserEmailUseCase {
	return &changeUserEmailUseCase{userRepository: userRepository}
}

func (u *changeUserEmailUseCase) Execute(ctx context.Context, id int64, email *domain.Email, userID int64) error {
	if err := domain.ValidateAssignedUserID(id); err != nil {
		return err
	}

	userFound, err := u.userRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingEmail, err)
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
			return fmt.Errorf("%v: %w", domain.ErrChangingEmail, err)
		}
		if foundEmail != nil {
			return domain.ErrUserEmailAlreadyExists
		}
	}

	updatedUser, err := userFound.ChangeEmail(email, userID)
	if err != nil {
		return err
	}

	err = u.userRepository.Update(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrChangingEmail, err)
	}

	return nil
}
