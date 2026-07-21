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
		name    string
		setup   func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository)
		input   RequestOTPInput
		want    *RequestOTPOutput
		wantErr error
	}{
		{
			name: "success - new OTP",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, nil)
				otpRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("SendOTP", mock.Anything, "+573001234567", mock.Anything).Return(nil)
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    &RequestOTPOutput{},
			wantErr: nil,
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
					false,
					time.Time{},
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(existingOTP, nil)
				otpRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("SendOTP", mock.Anything, "+573001234567", mock.Anything).Return(nil)
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    &RequestOTPOutput{},
			wantErr: nil,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(nil, nil)
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    nil,
			wantErr: ErrRequestOTP,
		},
		{
			name: "error - FindByPhone fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(nil, errors.New("database error"))
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    nil,
			wantErr: ErrRequestOTP,
		},
		{
			name: "error - Create fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, nil)
				otpRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    nil,
			wantErr: ErrRequestOTP,
		},
		{
			name: "error - SendOTP fails",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(nil, nil)
				otpRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				otpSender.On("SendOTP", mock.Anything, "+573001234567", mock.Anything).Return(errors.New("sms error"))
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    nil,
			wantErr: ErrRequestOTP,
		},
		{
			name: "error - invalid phone number",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				userRepo.On("FindByPhone", mock.Anything, "").Return(nil, nil)
			},
			input:   RequestOTPInput{PhoneNumber: ""},
			want:    nil,
			wantErr: ErrRequestOTP,
		},
		{
			name: "error - OTP blocked",
			setup: func(t *testing.T, userRepo *mock_domain.MockUserRepository, otpRepo *mock_domain.MockAuthOTPRepository, otpSender *mock_domain.MockOTPSenderRepository) {
				blockedOTP, _ := domain.NewAuthOTPWithSession(
					"blocked-session-id",
					"+573001234567",
					domain.OTPCode("123456"),
					time.Now().Add(-2*time.Minute),
					time.Now().Add(3*time.Minute),
					time.Time{},
					5,
					3,
					true,
					time.Now().Add(4*time.Hour),
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(blockedOTP, nil)
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    nil,
			wantErr: ErrRequestOTP,
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
					false,
					time.Time{},
				)
				userRepo.On("FindByPhone", mock.Anything, "+573001234567").Return(user, nil)
				otpRepo.On("FindByPhoneNumber", mock.Anything, "+573001234567").Return(maxedOTP, nil)
			},
			input:   RequestOTPInput{PhoneNumber: "+573001234567"},
			want:    nil,
			wantErr: ErrRequestOTP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mock_domain.NewMockUserRepository(t)
			otpRepo := mock_domain.NewMockAuthOTPRepository(t)
			otpSender := mock_domain.NewMockOTPSenderRepository(t)
			tt.setup(t, userRepo, otpRepo, otpSender)

			uc := NewRequestOTP(userRepo, otpRepo, otpSender)
			got, err := uc.Execute(context.Background(), tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && got == nil {
				t.Errorf("Execute() returned nil, want non-nil")
			}

			if tt.wantErr == nil && got.SessionID == "" {
				t.Errorf("Execute() SessionID should not be empty")
			}
		})
	}
}
