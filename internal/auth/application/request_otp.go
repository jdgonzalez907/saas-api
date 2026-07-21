package application

import (
	"context"
	"fmt"
	"time"

	"github.com/jdgonzalez907/saas-api/internal/auth/domain"
)

type RequestOTP interface {
	Execute(ctx context.Context, phoneNumber string) (*domain.AuthOTP, error)
}

type requestOTP struct {
	userRepo  domain.UserRepository
	otpRepo   domain.AuthOTPRepository
	otpSender domain.OTPSenderRepository
}

func NewRequestOTP(
	userRepo domain.UserRepository,
	otpRepo domain.AuthOTPRepository,
	otpSender domain.OTPSenderRepository,
) RequestOTP {
	return &requestOTP{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		otpSender: otpSender,
	}
}

func (uc *requestOTP) Execute(ctx context.Context, phoneNumber string) (*domain.AuthOTP, error) {
	user, err := uc.userRepo.FindByPhone(ctx, phoneNumber)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}

	existingOTP, err := uc.otpRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, uc.wrapError(err)
	}

	if existingOTP != nil {
		if err := existingOTP.Generate(time.Now()); err != nil {
			return nil, uc.wrapError(err)
		}
		if err := uc.otpRepo.Update(ctx, existingOTP); err != nil {
			return nil, uc.wrapError(err)
		}
		if err := uc.otpSender.Send(ctx, existingOTP); err != nil {
			return nil, uc.wrapError(err)
		}
		return existingOTP, nil
	}

	newOTP, err := domain.NewAuthOTP(phoneNumber)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if err := uc.otpRepo.Create(ctx, newOTP); err != nil {
		return nil, uc.wrapError(err)
	}
	if err := uc.otpSender.Send(ctx, newOTP); err != nil {
		return nil, uc.wrapError(err)
	}
	return newOTP, nil
}

func (uc *requestOTP) wrapError(err error) error {
	return fmt.Errorf("%w: %v", domain.ErrRequestOTP, err)
}
