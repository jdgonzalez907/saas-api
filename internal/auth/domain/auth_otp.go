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
)

const (
	OTPExpirationMinutes  = 5
	MaxResends            = 3
	BlockDurationHours    = 4
	ResendCooldownMinutes = 1
)

type AuthOTP struct {
	id             string
	phoneNumber    string
	code           OTPCode
	createdAt      time.Time
	expiresAt      time.Time
	lastResentAt   time.Time
	resendCount    int
	failedAttempts int
	isBlocked      bool
	blockedUntil   time.Time
}

type AuthOTPDTO struct {
	ID             string    `json:"id"`
	PhoneNumber    string    `json:"phone_number"`
	Code           string    `json:"code"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	LastResentAt   time.Time `json:"last_resent_at"`
	ResendCount    int       `json:"resend_count"`
	FailedAttempts int       `json:"failed_attempts"`
	IsBlocked      bool      `json:"is_blocked"`
	BlockedUntil   time.Time `json:"blocked_until"`
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
		id:          uuid.New().String(),
		phoneNumber: phoneNumber,
		code:        code,
		createdAt:   now,
		expiresAt:   now.Add(OTPExpirationMinutes * time.Minute),
	}, nil
}

func NewAuthOTPWithSession(
	id, phoneNumber string,
	code OTPCode,
	createdAt, expiresAt, lastResentAt time.Time,
	resendCount, failedAttempts int,
	isBlocked bool,
	blockedUntil time.Time,
) (*AuthOTP, error) {
	if id == "" {
		return nil, ErrSessionIDRequired
	}
	if phoneNumber == "" {
		return nil, ErrInvalidPhoneNumber
	}

	return &AuthOTP{
		id:             id,
		phoneNumber:    phoneNumber,
		code:           code,
		createdAt:      createdAt,
		expiresAt:      expiresAt,
		lastResentAt:   lastResentAt,
		resendCount:    resendCount,
		failedAttempts: failedAttempts,
		isBlocked:      isBlocked,
		blockedUntil:   blockedUntil,
	}, nil
}

func (a *AuthOTP) ID() string              { return a.id }
func (a *AuthOTP) PhoneNumber() string     { return a.phoneNumber }
func (a *AuthOTP) Code() OTPCode           { return a.code }
func (a *AuthOTP) CreatedAt() time.Time    { return a.createdAt }
func (a *AuthOTP) ExpiresAt() time.Time    { return a.expiresAt }
func (a *AuthOTP) LastResentAt() time.Time { return a.lastResentAt }
func (a *AuthOTP) ResendCount() int        { return a.resendCount }
func (a *AuthOTP) FailedAttempts() int     { return a.failedAttempts }
func (a *AuthOTP) IsBlocked() bool         { return a.isBlocked }
func (a *AuthOTP) BlockedUntil() time.Time { return a.blockedUntil }

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

	a.clearBlock()

	if a.isExpired(now) {
		code, err := NewOTPCodeRandom()
		if err != nil {
			return err
		}
		a.code = code
		a.resendCount = 0
		a.failedAttempts = 0
		a.expiresAt = now.Add(OTPExpirationMinutes * time.Minute)
		a.lastResentAt = time.Time{}
		return nil
	}

	if a.canResend(now) {
		code, err := NewOTPCodeRandom()
		if err != nil {
			return err
		}
		a.code = code
		a.resendCount++
		a.lastResentAt = now
		a.expiresAt = now.Add(OTPExpirationMinutes * time.Minute)
		return nil
	}

	return ErrOTPMaxResendsReached
}

func (a *AuthOTP) ToDTO() AuthOTPDTO {
	return AuthOTPDTO{
		ID:             a.id,
		PhoneNumber:    a.phoneNumber,
		Code:           a.code.Value(),
		CreatedAt:      a.createdAt,
		ExpiresAt:      a.expiresAt,
		LastResentAt:   a.lastResentAt,
		ResendCount:    a.resendCount,
		FailedAttempts: a.failedAttempts,
		IsBlocked:      a.isBlocked,
		BlockedUntil:   a.blockedUntil,
	}
}

func (a *AuthOTP) isBlockedNow(now time.Time) bool {
	return a.IsBlocked() && now.Before(a.BlockedUntil())
}

func (a *AuthOTP) clearBlock() {
	a.isBlocked = false
	a.blockedUntil = time.Time{}
	a.failedAttempts = 0
}

func (a *AuthOTP) isExpired(now time.Time) bool {
	return now.After(a.expiresAt)
}

func (a *AuthOTP) canResend(now time.Time) bool {
	if a.resendCount >= MaxResends {
		return false
	}
	hasNeverResent := a.lastResentAt.IsZero()
	isInCooldown := !hasNeverResent && now.Sub(a.lastResentAt) < ResendCooldownMinutes*time.Minute
	return !isInCooldown
}
