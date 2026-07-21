package domain

import "context"

type AuthOTPRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) (*AuthOTP, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*AuthOTP, error)
	Create(ctx context.Context, otp *AuthOTP) error
	Update(ctx context.Context, otp *AuthOTP) error
}
