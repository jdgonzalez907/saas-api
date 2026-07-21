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

func TestRequestOTP_Execute(t *testing.T) {
	user, _ := domain.NewUser(1, "John Doe")

	tests := []struct {
		name        string
		setup       func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository)
		phoneNumber string
		wantID      string
		wantErr     error
	}{
		{
			name: "success - new OTP",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, nil)
				otpRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("Send", mock.Anything, mock.Anything).Return(nil)
			},
			phoneNumber: "+573001234567",
			wantErr:     nil,
		},
		{
			name: "success - resend OTP",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				existingOTP, _ := domain.NewAuthOTPWithSession(
					"existing-session-id",
					"+573001234567",
					domain.OTPCode("123456"),
					time.Now().Add(-2*time.Minute),
					time.Now().Add(3*time.Minute),
					time.Now().Add(-2*time.Minute),
					1,
					0,
					nil,
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(existingOTP, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("Send", mock.Anything, mock.Anything).Return(nil)
			},
			phoneNumber: "+573001234567",
			wantID:      "existing-session-id",
			wantErr:     nil,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(nil, nil)
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - FindByPhone fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(nil, errors.New("database error"))
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - FindByPhoneNumber fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, errors.New("database error"))
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - Create fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, nil)
				otpRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - SendOTP fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, nil)
				otpRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("Send", mock.Anything, mock.Anything).Return(errors.New("sms error"))
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - invalid phone number",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "").Return(nil, nil)
			},
			phoneNumber: "",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - OTP blocked",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				blockedUntil := time.Now().Add(4 * time.Hour)
				blockedOTP, _ := domain.NewAuthOTPWithSession(
					"blocked-session-id",
					"+573001234567",
					domain.OTPCode("123456"),
					time.Now().Add(-2*time.Minute),
					time.Now().Add(3*time.Minute),
					time.Time{},
					5,
					3,
					&blockedUntil,
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(blockedOTP, nil)
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - max resends reached",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				maxedOTP, _ := domain.NewAuthOTPWithSession(
					"maxed-session-id",
					"+573001234567",
					domain.OTPCode("123456"),
					time.Now().Add(-2*time.Minute),
					time.Now().Add(3*time.Minute),
					time.Now().Add(-2*time.Minute),
					3,
					0,
					nil,
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(maxedOTP, nil)
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "error - Update fails on resend",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				existingOTP, _ := domain.NewAuthOTPWithSession(
					"existing-session-id",
					"+573001234567",
					domain.OTPCode("123456"),
					time.Now().Add(-2*time.Minute),
					time.Now().Add(3*time.Minute),
					time.Now().Add(-2*time.Minute),
					1,
					0,
					nil,
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(existingOTP, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			phoneNumber: "+573001234567",
			wantErr:     domain.ErrRequestOTP,
		},
		{
			name: "success - returns OTP entity",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				existingOTP, _ := domain.NewAuthOTPWithSession(
					"existing-session-id",
					"+573001234567",
					domain.OTPCode("123456"),
					time.Now().Add(-2*time.Minute),
					time.Now().Add(3*time.Minute),
					time.Now().Add(-2*time.Minute),
					1,
					0,
					nil,
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(existingOTP, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("Send", mock.Anything, mock.Anything).Return(nil)
			},
			phoneNumber: "+573001234567",
			wantID:      "existing-session-id",
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mock_domain.NewMockUserRepository(t)
			otpRepo := mock_domain.NewMockAuthOTPRepository(t)
			otpSender := mock_domain.NewMockOTPSenderRepository(t)
			tt.setup(t, userRepo, otpRepo, otpSender)

			uc := NewRequestOTP(userRepo, otpRepo, otpSender)
			got, err := uc.Execute(context.Background(), tt.phoneNumber)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && got == nil {
				t.Errorf("Execute() returned nil, want non-nil")
			}

			if tt.wantErr == nil && got.ID() == "" {
				t.Errorf("Execute() should return OTP with non-empty ID")
			}

			if tt.wantID != "" && got != nil && got.ID() != tt.wantID {
				t.Errorf("Execute() ID = %v, want %v", got.ID(), tt.wantID)
			}
		})
	}
}
