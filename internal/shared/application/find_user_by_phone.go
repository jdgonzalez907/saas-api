package application

import (
	"context"

	"jdgonzalez907/saas-api/internal/shared/domain"
	userDomain "jdgonzalez907/saas-api/internal/users/domain"
)

type FindUserByPhoneUseCase interface {
	Execute(ctx context.Context, countryCode, number string) (*domain.User, error)
}

type findUserByPhoneUseCase struct {
	userRepository userDomain.UserRepository
}

func NewFindUserByPhoneUseCase(userRepository userDomain.UserRepository) FindUserByPhoneUseCase {
	return &findUserByPhoneUseCase{userRepository: userRepository}
}

func (r *findUserByPhoneUseCase) Execute(ctx context.Context, countryCode, number string) (*domain.User, error) {
	phone, err := userDomain.NewPhone(countryCode, number)
	if err != nil {
		return nil, err
	}

	userFound, err := r.userRepository.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if userFound == nil {
		return nil, userDomain.ErrUserNotFound
	}

	var email *string
	if userFound.Email() != nil {
		emailVal := userFound.Email().Value()
		email = &emailVal
	}

	user := domain.NewUser(
		userFound.ID(),
		userFound.FullName(),
		userFound.Phone().CountryCode(),
		userFound.Phone().Number(),
		email,
	)

	return user, nil
}
