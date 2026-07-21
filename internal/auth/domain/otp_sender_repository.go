package domain

import "context"

type OTPSenderRepository interface {
	Send(ctx context.Context, otp *AuthOTP) error
}
