package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOTPExpired           = errors.New("OTP has expired")
	ErrOTPInvalid           = errors.New("OTP code is invalid")
	ErrOTPMaxResendsReached = errors.New("maximum OTP resend attempts reached")
	ErrOTPBlocked           = errors.New("OTP session is blocked")
	ErrInvalidPhoneNumber   = errors.New("invalid phone number format")
	ErrSessionIDRequired    = errors.New("session ID is required")
	ErrSessionNotFound      = errors.New("session not found")

	ErrRequestOTP = errors.New("cannot request OTP")
	ErrVerifyOTP  = errors.New("cannot verify OTP")
)

const (
	OTPExpirationMinutes  = 5
	MaxResends            = 3
	BlockDurationHours    = 4
	ResendCooldownMinutes = 1
	MaxFailedAttempts     = 3
)

type AuthOTP struct {
	id              string
	phoneNumber     string
	code            OTPCode
	createdAt       time.Time
	expiresAt       time.Time
	lastGeneratedAt time.Time
	resendCount     int
	failedAttempts  int
	blockedUntil    *time.Time
}

type AuthOTPDTO struct {
	ID              string     `json:"id"`
	PhoneNumber     string     `json:"phone_number"`
	Code            string     `json:"code"`
	CreatedAt       time.Time  `json:"created_at"`
	ExpiresAt       time.Time  `json:"expires_at"`
	LastGeneratedAt time.Time  `json:"last_generated_at"`
	ResendCount     int        `json:"resend_count"`
	FailedAttempts  int        `json:"failed_attempts"`
	BlockedUntil    *time.Time `json:"blocked_until"`
}

func NewAuthOTP(phoneNumber string) (*AuthOTP, error) {
	if phoneNumber == "" {
		return nil, ErrInvalidPhoneNumber
	}

	now := time.Now()
	code, err := NewOTPCodeRandom()
	if err != nil {
		return nil, err
	}

	return &AuthOTP{
		id:              uuid.New().String(),
		phoneNumber:     phoneNumber,
		code:            code,
		createdAt:       now,
		expiresAt:       now.Add(OTPExpirationMinutes * time.Minute),
		lastGeneratedAt: now,
		resendCount:     0,
		failedAttempts:  0,
		blockedUntil:    nil,
	}, nil
}

func NewAuthOTPWithSession(
	id, phoneNumber string,
	code OTPCode,
	createdAt, expiresAt, lastGeneratedAt time.Time,
	resendCount, failedAttempts int,
	blockedUntil *time.Time,
) (*AuthOTP, error) {
	if id == "" {
		return nil, ErrSessionIDRequired
	}
	if phoneNumber == "" {
		return nil, ErrInvalidPhoneNumber
	}

	return &AuthOTP{
		id:              id,
		phoneNumber:     phoneNumber,
		code:            code,
		createdAt:       createdAt,
		expiresAt:       expiresAt,
		lastGeneratedAt: lastGeneratedAt,
		resendCount:     resendCount,
		failedAttempts:  failedAttempts,
		blockedUntil:    blockedUntil,
	}, nil
}

func (a *AuthOTP) ID() string                 { return a.id }
func (a *AuthOTP) PhoneNumber() string        { return a.phoneNumber }
func (a *AuthOTP) Code() OTPCode              { return a.code }
func (a *AuthOTP) CreatedAt() time.Time       { return a.createdAt }
func (a *AuthOTP) ExpiresAt() time.Time       { return a.expiresAt }
func (a *AuthOTP) LastGeneratedAt() time.Time { return a.lastGeneratedAt }
func (a *AuthOTP) ResendCount() int           { return a.resendCount }
func (a *AuthOTP) FailedAttempts() int        { return a.failedAttempts }
func (a *AuthOTP) BlockedUntil() *time.Time   { return a.blockedUntil }

func (a *AuthOTP) Equals(other *AuthOTP) bool {
	if other == nil {
		return false
	}
	return a.id == other.id
}

func (a *AuthOTP) Generate(now time.Time) error {
	if a.isBlockedNow(now) {
		return ErrOTPBlocked
	}

	if a.isExpired(now) {
		code, err := NewOTPCodeRandom()
		if err != nil {
			return err
		}
		a.code = code
		a.createdAt = now
		a.expiresAt = now.Add(OTPExpirationMinutes * time.Minute)
		a.lastGeneratedAt = now
		a.resendCount = 0
		a.failedAttempts = 0
		a.blockedUntil = nil

		return nil
	}

	if a.canResend(now) {
		code, err := NewOTPCodeRandom()
		if err != nil {
			return err
		}
		a.code = code
		a.expiresAt = now.Add(OTPExpirationMinutes * time.Minute)
		a.lastGeneratedAt = now
		a.resendCount++
		a.failedAttempts = 0
		a.blockedUntil = nil

		return nil
	}

	a.block(now)
	return ErrOTPMaxResendsReached
}

func (a *AuthOTP) ToDTO() AuthOTPDTO {
	return AuthOTPDTO{
		ID:              a.id,
		PhoneNumber:     a.phoneNumber,
		Code:            a.code.Value(),
		CreatedAt:       a.createdAt,
		ExpiresAt:       a.expiresAt,
		LastGeneratedAt: a.lastGeneratedAt,
		ResendCount:     a.resendCount,
		FailedAttempts:  a.failedAttempts,
		BlockedUntil:    a.blockedUntil,
	}
}

func (a *AuthOTP) Verify(code string, now time.Time) error {
	if a.isBlockedNow(now) {
		return ErrOTPBlocked
	}
	if a.isExpired(now) {
		return ErrOTPExpired
	}
	otpCode, err := NewOTPCode(code)
	if err != nil {
		a.fail(now)
		return ErrOTPInvalid
	}
	if !a.code.Equals(otpCode) {
		a.fail(now)
		return ErrOTPInvalid
	}
	a.invalidate(now)
	return nil
}

func (a *AuthOTP) isBlockedNow(now time.Time) bool {
	return a.blockedUntil != nil && now.Before(*a.blockedUntil)
}

func (a *AuthOTP) block(now time.Time) {
	t := now.Add(BlockDurationHours * time.Hour)
	a.blockedUntil = &t
}

func (a *AuthOTP) fail(now time.Time) {
	a.failedAttempts++
	if a.failedAttempts >= MaxFailedAttempts {
		a.block(now)
	}
}

func (a *AuthOTP) invalidate(now time.Time) {
	a.expiresAt = now.Add(-1 * time.Nanosecond)
}

func (a *AuthOTP) isExpired(now time.Time) bool {
	return now.After(a.expiresAt)
}

func (a *AuthOTP) canResend(now time.Time) bool {
	if a.resendCount >= MaxResends {
		return false
	}
	isInCooldown := now.Sub(a.lastGeneratedAt) < ResendCooldownMinutes*time.Minute
	return !isInCooldown
}
