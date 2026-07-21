package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jdgonzalez907/saas-api/internal/auth/domain"
	mock_domain "github.com/jdgonzalez907/saas-api/mocks/auth/domain"
	"github.com/stretchr/testify/mock"
)

func TestVerifyOTP_Execute(t *testing.T) {
	user, _ := domain.NewUser(1, "John Doe")
	code, _ := domain.NewOTPCode("123456")

	tests := []struct {
		name      string
		setup     func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository)
		sessionID string
		code      string
		wantErr   error
	}{
		{
			name: "success - verifies OTP and returns user",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   nil,
		},
		{
			name: "error - session not found",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otpRepo.On("FindBySessionID", mock.Anything, "nonexistent-session").Return(nil, nil)
			},
			sessionID: "nonexistent-session",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(nil, nil)
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - OTP expired",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now().Add(-10*time.Minute), time.Now().Add(-5*time.Minute), time.Now().Add(-5*time.Minute), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - OTP blocked",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				blockedUntil := time.Now().Add(4 * time.Hour)
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, &blockedUntil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - invalid code",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			sessionID: "test-session-id",
			code:      "999999",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - FindBySessionID fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(nil, errors.New("database error"))
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - FindByPhone fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(nil, errors.New("database error"))
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - Update fails on valid code",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			sessionID: "test-session-id",
			code:      "123456",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - Update fails on invalid code",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			sessionID: "test-session-id",
			code:      "999999",
			wantErr:   domain.ErrVerifyOTP,
		},
		{
			name: "error - invalid code format",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository) {
				otp := mustNewAuthOTPWithSession(t, "test-session-id", "+573001234567", code, time.Now(), time.Now().Add(5*time.Minute), time.Now(), 0, 0, nil)
				otpRepo.On("FindBySessionID", mock.Anything, "test-session-id").Return(otp, nil)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			sessionID: "test-session-id",
			code:      "12345",
			wantErr:   domain.ErrVerifyOTP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mock_domain.NewMockUserRepository(t)
			otpRepo := mock_domain.NewMockAuthOTPRepository(t)
			tt.setup(t, userRepo, otpRepo)

			uc := NewVerifyOTP(userRepo, otpRepo)
			got, err := uc.Execute(context.Background(), tt.sessionID, tt.code)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && got == nil {
				t.Errorf("Execute() returned nil, want non-nil")
			}

			if tt.wantErr == nil && got.ID() != 1 {
				t.Errorf("Execute() returned user with ID = %v, want 1", got.ID())
			}
		})
	}
}

func mustNewAuthOTPWithSession(
	t *testing.T,
	id, phoneNumber string,
	code domain.OTPCode,
	createdAt, expiresAt, lastGeneratedAt time.Time,
	resendCount, failedAttempts int,
	blockedUntil *time.Time,
) *domain.AuthOTP {
	t.Helper()
	otp, err := domain.NewAuthOTPWithSession(
		id,
		phoneNumber,
		code,
		createdAt,
		expiresAt,
		lastGeneratedAt,
		resendCount,
		failedAttempts,
		blockedUntil,
	)
	if err != nil {
		t.Fatalf("mustNewAuthOTPWithSession() error = %v", err)
	}
	return otp
}
