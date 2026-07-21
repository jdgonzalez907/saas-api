package application

import (
	"context"
	"fmt"
	"time"

	"github.com/jdgonzalez907/saas-api/internal/auth/domain"
)

type VerifyOTP interface {
	Execute(ctx context.Context, sessionID string, code string) (*domain.User, error)
}

type verifyOTP struct {
	userRepo domain.UserRepository
	otpRepo  domain.AuthOTPRepository
}

func NewVerifyOTP(
	userRepo domain.UserRepository,
	otpRepo domain.AuthOTPRepository,
) VerifyOTP {
	return &verifyOTP{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (uc *verifyOTP) Execute(ctx context.Context, sessionID string, code string) (*domain.User, error) {
	otp, err := uc.otpRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if otp == nil {
		return nil, uc.wrapError(domain.ErrSessionNotFound)
	}

	user, err := uc.userRepo.FindByPhone(ctx, otp.PhoneNumber())
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}

	err = otp.Verify(code, time.Now())
	if updateErr := uc.otpRepo.Update(ctx, otp); updateErr != nil {
		return nil, uc.wrapError(updateErr)
	}
	if err != nil {
		return nil, uc.wrapError(err)
	}

	return user, nil
}

func (uc *verifyOTP) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrVerifyOTP, err)
}
