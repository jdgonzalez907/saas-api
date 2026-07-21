package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrOTPCodeEmpty        = errors.New("OTP code cannot be empty")
	ErrOTPCodeInvalid      = errors.New("OTP code must be exactly 6 digits")
	ErrOTPGenerationFailed = errors.New("failed to generate OTP code")
)

const OTPCodeLength = 6

type OTPCode string

func NewOTPCode(code string) (OTPCode, error) {
	if code == "" {
		return "", ErrOTPCodeEmpty
	}

	if len(code) != OTPCodeLength {
		return "", ErrOTPCodeInvalid
	}

	for i := 0; i < len(code); i++ {
		if code[i] < '0' || code[i] > '9' {
			return "", ErrOTPCodeInvalid
		}
	}

	return OTPCode(code), nil
}

func NewOTPCodeRandom() (OTPCode, error) {
	code := make([]byte, OTPCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrOTPGenerationFailed, err)
		}
		code[i] = byte('0' + n.Int64())
	}
	return OTPCode(string(code)), nil
}

func NewOTPCodeFromString(s string) OTPCode {
	return OTPCode(s)
}

func (v OTPCode) Value() string {
	return string(v)
}

func (v OTPCode) Equals(other OTPCode) bool {
	return v == other
}

func (v OTPCode) ToDTO() string {
	return string(v)
}
