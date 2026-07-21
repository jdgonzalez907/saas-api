package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jdgonzalez907/saas-api/internal/auth/domain"
)

var ErrRequestOTP = errors.New("cannot request OTP")

type RequestOTP interface {
	Execute(ctx context.Context, input RequestOTPInput) (*RequestOTPOutput, error)
}

type RequestOTPInput struct {
	PhoneNumber string
}

type RequestOTPOutput struct {
	SessionID string
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

func (uc *requestOTP) Execute(ctx context.Context, input RequestOTPInput) (*RequestOTPOutput, error) {
	user, err := uc.userRepo.FindByPhone(ctx, input.PhoneNumber)
	if err != nil {
		return nil, uc.wrapError(err)
	}
	if user == nil {
		return nil, uc.wrapError(domain.ErrUserNotFound)
	}

	existingOTP, err := uc.otpRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, uc.wrapError(err)
	}

	var otp *domain.AuthOTP
	if existingOTP != nil {
		otp = existingOTP
	} else {
		otp, err = domain.NewAuthOTP(input.PhoneNumber)
		if err != nil {
			return nil, uc.wrapError(err)
		}
	}

	if err := otp.Generate(time.Now()); err != nil {
		return nil, uc.wrapError(err)
	}

	if existingOTP != nil {
		err = uc.otpRepo.Update(ctx, otp)
	} else {
		err = uc.otpRepo.Create(ctx, otp)
	}
	if err != nil {
		return nil, uc.wrapError(err)
	}

	if err := uc.otpSender.SendOTP(ctx, input.PhoneNumber, otp.Code().Value()); err != nil {
		return nil, uc.wrapError(err)
	}

	return &RequestOTPOutput{SessionID: otp.ID()}, nil
}

func (uc *requestOTP) wrapError(err error) error {
	return fmt.Errorf("%w: %v", ErrRequestOTP, err)
}
