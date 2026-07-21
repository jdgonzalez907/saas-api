package domain

import (
	"testing"
	"time"
)

func TestNewAuthOTP(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		wantErr     error
	}{
		{
			name:        "success - valid phone number",
			phoneNumber: "+573001234567",
			wantErr:     nil,
		},
		{
			name:        "success - US phone number",
			phoneNumber: "+12025551234",
			wantErr:     nil,
		},
		{
			name:        "error - empty phone number",
			phoneNumber: "",
			wantErr:     ErrInvalidPhoneNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := NewAuthOTP(tt.phoneNumber)
			if err != tt.wantErr {
				t.Errorf("NewAuthOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if otp.ID() == "" {
					t.Errorf("NewAuthOTP().ID() should not be empty")
				}
				if otp.Code() == "" {
					t.Errorf("NewAuthOTP().Code() should not be empty")
				}
			}
		})
	}
}

func TestNewAuthOTPWithSession(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	tests := []struct {
		name           string
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
		wantErr        error
	}{
		{
			name:           "success",
			id:             "test-session-id",
			phoneNumber:    "+573001234567",
			code:           code,
			createdAt:      now,
			expiresAt:      now.Add(5 * time.Minute),
			lastResentAt:   time.Time{},
			resendCount:    1,
			failedAttempts: 0,
			isBlocked:      false,
			blockedUntil:   time.Time{},
			wantErr:        nil,
		},
		{
			name:           "success - with failed attempts",
			id:             "test-session-id-2",
			phoneNumber:    "+573001234567",
			code:           code,
			createdAt:      now,
			expiresAt:      now.Add(5 * time.Minute),
			lastResentAt:   now.Add(-2 * time.Minute),
			resendCount:    2,
			failedAttempts: 1,
			isBlocked:      false,
			blockedUntil:   time.Time{},
			wantErr:        nil,
		},
		{
			name:           "success - blocked session",
			id:             "test-session-id-3",
			phoneNumber:    "+573001234567",
			code:           code,
			createdAt:      now,
			expiresAt:      now.Add(5 * time.Minute),
			lastResentAt:   time.Time{},
			resendCount:    5,
			failedAttempts: 3,
			isBlocked:      true,
			blockedUntil:   now.Add(4 * time.Hour),
			wantErr:        nil,
		},
		{
			name:           "error - empty ID",
			id:             "",
			phoneNumber:    "+573001234567",
			code:           code,
			createdAt:      now,
			expiresAt:      now.Add(5 * time.Minute),
			lastResentAt:   time.Time{},
			resendCount:    0,
			failedAttempts: 0,
			isBlocked:      false,
			blockedUntil:   time.Time{},
			wantErr:        ErrSessionIDRequired,
		},
		{
			name:           "error - empty phone number",
			id:             "test-session-id-4",
			phoneNumber:    "",
			code:           code,
			createdAt:      now,
			expiresAt:      now.Add(5 * time.Minute),
			lastResentAt:   time.Time{},
			resendCount:    0,
			failedAttempts: 0,
			isBlocked:      false,
			blockedUntil:   time.Time{},
			wantErr:        ErrInvalidPhoneNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := NewAuthOTPWithSession(
				tt.id,
				tt.phoneNumber,
				tt.code,
				tt.createdAt,
				tt.expiresAt,
				tt.lastResentAt,
				tt.resendCount,
				tt.failedAttempts,
				tt.isBlocked,
				tt.blockedUntil,
			)
			if err != tt.wantErr {
				t.Errorf("NewAuthOTPWithSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if otp.ID() != tt.id {
					t.Errorf("NewAuthOTPWithSession().ID() = %v, want %v", otp.ID(), tt.id)
				}
				if otp.Code() != tt.code {
					t.Errorf("NewAuthOTPWithSession().Code() = %v, want %v", otp.Code(), tt.code)
				}
			}
		})
	}
}

func TestAuthOTP_Equals(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	tests := []struct {
		name string
		e1   *AuthOTP
		e2   *AuthOTP
		want bool
	}{
		{
			name: "equal - same ID",
			e1:   mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), time.Time{}, 0, 0, false, time.Time{}),
			e2:   mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), time.Time{}, 0, 0, false, time.Time{}),
			want: true,
		},
		{
			name: "not equal - different ID",
			e1:   mustNewAuthOTPWithSession(t, "test-id-1", "+573001234567", code, now, now.Add(5*time.Minute), time.Time{}, 0, 0, false, time.Time{}),
			e2:   mustNewAuthOTPWithSession(t, "test-id-2", "+573001234567", code, now, now.Add(5*time.Minute), time.Time{}, 0, 0, false, time.Time{}),
			want: false,
		},
		{
			name: "not equal - nil other",
			e1:   mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), time.Time{}, 0, 0, false, time.Time{}),
			e2:   nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e1.Equals(tt.e2); got != tt.want {
				t.Errorf("AuthOTP.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthOTP_Generate(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	tests := []struct {
		name    string
		otp     *AuthOTP
		now     time.Time
		wantErr error
	}{
		{
			name:    "success - generate new code for expired OTP",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-10*time.Minute), now.Add(-5*time.Minute), time.Time{}, 0, 0, false, time.Time{}),
			now:     now,
			wantErr: nil,
		},
		{
			name:    "success - generate new code for resend",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-2*time.Minute), 1, 0, false, time.Time{}),
			now:     now,
			wantErr: nil,
		},
		{
			name:    "error - blocked session",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), time.Time{}, 5, 3, true, now.Add(4*time.Hour)),
			now:     now,
			wantErr: ErrOTPBlocked,
		},
		{
			name:    "error - max resends reached",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-2*time.Minute), 3, 0, false, time.Time{}),
			now:     now,
			wantErr: ErrOTPMaxResendsReached,
		},
		{
			name:    "success - clear block after expiration",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-10*time.Minute), now.Add(-5*time.Minute), time.Time{}, 5, 3, true, now.Add(-1*time.Hour)),
			now:     now,
			wantErr: nil,
		},
		{
			name:    "error - resend cooldown not elapsed",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-30*time.Second), 1, 0, false, time.Time{}),
			now:     now,
			wantErr: ErrOTPMaxResendsReached,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCode := tt.otp.Code()
			err := tt.otp.Generate(tt.now)
			if err != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if tt.otp.Code() == originalCode {
					t.Errorf("Generate() should change the code")
				}
			}
		})
	}
}

func TestAuthOTP_ToDTO(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")
	lastResent := now.Add(-2 * time.Minute)
	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), lastResent, 1, 2, true, now.Add(4*time.Hour))

	dto := otp.ToDTO()

	if dto.ID != "test-id" {
		t.Errorf("ToDTO().ID = %v, want %v", dto.ID, "test-id")
	}
	if dto.PhoneNumber != "+573001234567" {
		t.Errorf("ToDTO().PhoneNumber = %v, want %v", dto.PhoneNumber, "+573001234567")
	}
	if dto.Code != "123456" {
		t.Errorf("ToDTO().Code = %v, want %v", dto.Code, "123456")
	}
	if !dto.CreatedAt.Equal(now) {
		t.Errorf("ToDTO().CreatedAt = %v, want %v", dto.CreatedAt, now)
	}
	if !dto.ExpiresAt.Equal(now.Add(5 * time.Minute)) {
		t.Errorf("ToDTO().ExpiresAt = %v, want %v", dto.ExpiresAt, now.Add(5*time.Minute))
	}
	if !dto.LastResentAt.Equal(lastResent) {
		t.Errorf("ToDTO().LastResentAt = %v, want %v", dto.LastResentAt, lastResent)
	}
	if dto.ResendCount != 1 {
		t.Errorf("ToDTO().ResendCount = %v, want %v", dto.ResendCount, 1)
	}
	if dto.FailedAttempts != 2 {
		t.Errorf("ToDTO().FailedAttempts = %v, want %v", dto.FailedAttempts, 2)
	}
	if !dto.IsBlocked {
		t.Errorf("ToDTO().IsBlocked should be true")
	}
	if !dto.BlockedUntil.Equal(now.Add(4 * time.Hour)) {
		t.Errorf("ToDTO().BlockedUntil = %v, want %v", dto.BlockedUntil, now.Add(4*time.Hour))
	}
}

func mustNewAuthOTPWithSession(
	t *testing.T,
	id, phoneNumber string,
	code OTPCode,
	createdAt, expiresAt, lastResentAt time.Time,
	resendCount, failedAttempts int,
	isBlocked bool,
	blockedUntil time.Time,
) *AuthOTP {
	t.Helper()
	otp, err := NewAuthOTPWithSession(
		id,
		phoneNumber,
		code,
		createdAt,
		expiresAt,
		lastResentAt,
		resendCount,
		failedAttempts,
		isBlocked,
		blockedUntil,
	)
	if err != nil {
		t.Fatalf("mustNewAuthOTPWithSession() error = %v", err)
	}
	return otp
}
