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
				if otp.BlockedUntil() != nil {
					t.Errorf("NewAuthOTP().BlockedUntil() should be nil")
				}
			}
		})
	}
}

func TestNewAuthOTPWithSession(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	tests := []struct {
		name             string
		id               string
		phoneNumber      string
		code             OTPCode
		createdAt        time.Time
		expiresAt        time.Time
		lastGeneratedAt  time.Time
		resendCount      int
		failedAttempts   int
		blockedUntil     *time.Time
		wantErr          error
	}{
		{
			name:            "success",
			id:              "test-session-id",
			phoneNumber:     "+573001234567",
			code:            code,
			createdAt:       now,
			expiresAt:       now.Add(5 * time.Minute),
			lastGeneratedAt: now,
			resendCount:     1,
			failedAttempts:  0,
			blockedUntil:    nil,
			wantErr:         nil,
		},
		{
			name:            "success - with failed attempts",
			id:              "test-session-id-2",
			phoneNumber:     "+573001234567",
			code:            code,
			createdAt:       now,
			expiresAt:       now.Add(5 * time.Minute),
			lastGeneratedAt: now.Add(-2 * time.Minute),
			resendCount:     2,
			failedAttempts:  1,
			blockedUntil:    nil,
			wantErr:         nil,
		},
		{
			name:            "success - blocked session",
			id:              "test-session-id-3",
			phoneNumber:     "+573001234567",
			code:            code,
			createdAt:       now,
			expiresAt:       now.Add(5 * time.Minute),
			lastGeneratedAt: time.Time{},
			resendCount:     5,
			failedAttempts:  3,
			blockedUntil:    ptrTime(now.Add(4 * time.Hour)),
			wantErr:         nil,
		},
		{
			name:            "error - empty ID",
			id:              "",
			phoneNumber:     "+573001234567",
			code:            code,
			createdAt:       now,
			expiresAt:       now.Add(5 * time.Minute),
			lastGeneratedAt: time.Time{},
			resendCount:     0,
			failedAttempts:  0,
			blockedUntil:    nil,
			wantErr:         ErrSessionIDRequired,
		},
		{
			name:            "error - empty phone number",
			id:              "test-session-id-4",
			phoneNumber:     "",
			code:            code,
			createdAt:       now,
			expiresAt:       now.Add(5 * time.Minute),
			lastGeneratedAt: time.Time{},
			resendCount:     0,
			failedAttempts:  0,
			blockedUntil:    nil,
			wantErr:         ErrInvalidPhoneNumber,
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
				tt.lastGeneratedAt,
				tt.resendCount,
				tt.failedAttempts,
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
			e1:   mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			e2:   mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			want: true,
		},
		{
			name: "not equal - different ID",
			e1:   mustNewAuthOTPWithSession(t, "test-id-1", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			e2:   mustNewAuthOTPWithSession(t, "test-id-2", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			want: false,
		},
		{
			name: "not equal - nil other",
			e1:   mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
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
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-10*time.Minute), now.Add(-5*time.Minute), now.Add(-5*time.Minute), 0, 0, nil),
			now:     now,
			wantErr: nil,
		},
		{
			name:    "success - generate new code for resend",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-2*time.Minute), 1, 0, nil),
			now:     now,
			wantErr: nil,
		},
		{
			name:    "error - blocked session",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), time.Time{}, 5, 3, ptrTime(now.Add(4*time.Hour))),
			now:     now,
			wantErr: ErrOTPBlocked,
		},
		{
			name:    "error - max resends reached",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-2*time.Minute), 3, 0, nil),
			now:     now,
			wantErr: ErrOTPMaxResendsReached,
		},
		{
			name:    "success - clear block after expiration",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-10*time.Minute), now.Add(-5*time.Minute), time.Time{}, 5, 3, ptrTime(now.Add(-1*time.Hour))),
			now:     now,
			wantErr: nil,
		},
		{
			name:    "error - resend cooldown not elapsed",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-30*time.Second), 1, 0, nil),
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
			if tt.name == "error - max resends reached" {
				if tt.otp.BlockedUntil() == nil {
					t.Error("Generate() should block session when max resends reached")
				}
			}
			if tt.name == "error - resend cooldown not elapsed" {
				if tt.otp.BlockedUntil() == nil {
					t.Error("Generate() should block session when cooldown not elapsed")
				}
			}
		})
	}
}

func TestAuthOTP_Generate_ExpiredResetsCounters(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code,
		now.Add(-10*time.Minute), now.Add(-5*time.Minute), now.Add(-5*time.Minute),
		3, 2, ptrTime(now.Add(-1*time.Hour)))

	err := otp.Generate(now)
	if err != nil {
		t.Fatalf("Generate() unexpected error = %v", err)
	}

	if otp.ResendCount() != 0 {
		t.Errorf("Generate() should reset ResendCount to 0, got %d", otp.ResendCount())
	}
	if otp.FailedAttempts() != 0 {
		t.Errorf("Generate() should reset FailedAttempts to 0, got %d", otp.FailedAttempts())
	}
	if otp.BlockedUntil() != nil {
		t.Errorf("Generate() should reset BlockedUntil to nil")
	}
}

func TestAuthOTP_Generate_ResendIncrementsCount(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code,
		now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-2*time.Minute),
		0, 0, nil)

	err := otp.Generate(now)
	if err != nil {
		t.Fatalf("Generate() unexpected error = %v", err)
	}

	if otp.ResendCount() != 1 {
		t.Errorf("Generate() should increment ResendCount to 1, got %d", otp.ResendCount())
	}
}

func TestAuthOTP_Generate_ExpiredDoesNotIncrementCount(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code,
		now.Add(-10*time.Minute), now.Add(-5*time.Minute), now.Add(-5*time.Minute),
		2, 0, nil)

	err := otp.Generate(now)
	if err != nil {
		t.Fatalf("Generate() unexpected error = %v", err)
	}

	if otp.ResendCount() != 0 {
		t.Errorf("Generate() on expired OTP should reset ResendCount to 0, got %d", otp.ResendCount())
	}
}

func TestAuthOTP_Getters(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")
	lastGenerated := now.Add(-2 * time.Minute)
	blockedUntil := now.Add(4 * time.Hour)
	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), lastGenerated, 1, 2, &blockedUntil)

	if otp.PhoneNumber() != "+573001234567" {
		t.Errorf("PhoneNumber() = %v, want %v", otp.PhoneNumber(), "+573001234567")
	}
	if !otp.CreatedAt().Equal(now) {
		t.Errorf("CreatedAt() = %v, want %v", otp.CreatedAt(), now)
	}
	if !otp.ExpiresAt().Equal(now.Add(5 * time.Minute)) {
		t.Errorf("ExpiresAt() = %v, want %v", otp.ExpiresAt(), now.Add(5*time.Minute))
	}
	if !otp.LastGeneratedAt().Equal(lastGenerated) {
		t.Errorf("LastGeneratedAt() = %v, want %v", otp.LastGeneratedAt(), lastGenerated)
	}
	if otp.ResendCount() != 1 {
		t.Errorf("ResendCount() = %v, want %v", otp.ResendCount(), 1)
	}
	if otp.FailedAttempts() != 2 {
		t.Errorf("FailedAttempts() = %v, want %v", otp.FailedAttempts(), 2)
	}
	if otp.BlockedUntil() == nil || !otp.BlockedUntil().Equal(blockedUntil) {
		t.Errorf("BlockedUntil() = %v, want %v", otp.BlockedUntil(), blockedUntil)
	}
}

func TestAuthOTP_Block(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")
	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil)

	otp.block(now)

	if otp.BlockedUntil() == nil {
		t.Fatal("block() should set BlockedUntil")
	}
	expected := now.Add(BlockDurationHours * time.Hour)
	if !otp.BlockedUntil().Equal(expected) {
		t.Errorf("block() BlockedUntil = %v, want %v", otp.BlockedUntil(), expected)
	}
}

func TestAuthOTP_ToDTO(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")
	lastGenerated := now.Add(-2 * time.Minute)
	blockedUntil := now.Add(4 * time.Hour)
	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), lastGenerated, 1, 2, &blockedUntil)

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
	if !dto.LastGeneratedAt.Equal(lastGenerated) {
		t.Errorf("ToDTO().LastGeneratedAt = %v, want %v", dto.LastGeneratedAt, lastGenerated)
	}
	if dto.ResendCount != 1 {
		t.Errorf("ToDTO().ResendCount = %v, want %v", dto.ResendCount, 1)
	}
	if dto.FailedAttempts != 2 {
		t.Errorf("ToDTO().FailedAttempts = %v, want %v", dto.FailedAttempts, 2)
	}
	if dto.BlockedUntil == nil || !dto.BlockedUntil.Equal(blockedUntil) {
		t.Errorf("ToDTO().BlockedUntil = %v, want %v", dto.BlockedUntil, blockedUntil)
	}
}

func TestAuthOTP_Generate_FirstTime(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	otp := mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code,
		now.Add(-2*time.Minute), now.Add(3*time.Minute), now.Add(-2*time.Minute),
		0, 0, nil)

	err := otp.Generate(now)
	if err != nil {
		t.Fatalf("Generate() unexpected error = %v", err)
	}

	if otp.ResendCount() != 1 {
		t.Errorf("Generate() should increment ResendCount to 1, got %d", otp.ResendCount())
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		fullname string
		wantErr  error
	}{
		{name: "success", id: 1, fullname: "John Doe", wantErr: nil},
		{name: "error - invalid ID", id: 0, fullname: "John Doe", wantErr: ErrUserInvalidID},
		{name: "error - negative ID", id: -1, fullname: "John Doe", wantErr: ErrUserInvalidID},
		{name: "error - empty name", id: 1, fullname: "", wantErr: ErrUserEmptyName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.fullname)
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if user.ID() != tt.id {
					t.Errorf("NewUser().ID() = %v, want %v", user.ID(), tt.id)
				}
				if user.Fullname() != tt.fullname {
					t.Errorf("NewUser().Fullname() = %v, want %v", user.Fullname(), tt.fullname)
				}
			}
		})
	}
}

func TestUser_Equals(t *testing.T) {
	user1, _ := NewUser(1, "John Doe")
	user2, _ := NewUser(1, "Jane Doe")
	user3, _ := NewUser(2, "John Doe")

	if !user1.Equals(user2) {
		t.Error("Equals() should return true for same ID")
	}
	if user1.Equals(user3) {
		t.Error("Equals() should return false for different ID")
	}
	if user1.Equals(nil) {
		t.Error("Equals() should return false for nil")
	}
}

func TestUser_ToDTO(t *testing.T) {
	user, _ := NewUser(1, "John Doe")
	dto := user.ToDTO()
	if dto.ID != 1 {
		t.Errorf("ToDTO().ID = %v, want %v", dto.ID, 1)
	}
	if dto.Fullname != "John Doe" {
		t.Errorf("ToDTO().Fullname = %v, want %v", dto.Fullname, "John Doe")
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func mustNewAuthOTPWithSession(
	t *testing.T,
	id, phoneNumber string,
	code OTPCode,
	createdAt, expiresAt, lastGeneratedAt time.Time,
	resendCount, failedAttempts int,
	blockedUntil *time.Time,
) *AuthOTP {
	t.Helper()
	otp, err := NewAuthOTPWithSession(
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

func TestAuthOTP_Verify(t *testing.T) {
	now := time.Now()
	code, _ := NewOTPCode("123456")

	tests := []struct {
		name    string
		otp     *AuthOTP
		code    string
		now     time.Time
		wantErr error
	}{
		{
			name:    "success - valid code",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			code:    "123456",
			now:     now,
			wantErr: nil,
		},
		{
			name:    "error - OTP blocked",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, ptrTime(now.Add(4*time.Hour))),
			code:    "123456",
			now:     now,
			wantErr: ErrOTPBlocked,
		},
		{
			name:    "error - OTP expired",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now.Add(-10*time.Minute), now.Add(-5*time.Minute), now.Add(-5*time.Minute), 0, 0, nil),
			code:    "123456",
			now:     now,
			wantErr: ErrOTPExpired,
		},
		{
			name:    "error - invalid code",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			code:    "999999",
			now:     now,
			wantErr: ErrOTPInvalid,
		},
		{
			name:    "error - invalid code blocks after max attempts",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 2, nil),
			code:    "999999",
			now:     now,
			wantErr: ErrOTPInvalid,
		},
		{
			name:    "error - invalid code format",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 0, nil),
			code:    "12345",
			now:     now,
			wantErr: ErrOTPInvalid,
		},
		{
			name:    "error - invalid code format blocks after max attempts",
			otp:     mustNewAuthOTPWithSession(t, "test-id", "+573001234567", code, now, now.Add(5*time.Minute), now, 0, 2, nil),
			code:    "12345",
			now:     now,
			wantErr: ErrOTPInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.otp.Verify(tt.code, tt.now)
			if err != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "success - valid code" {
				if !tt.otp.isExpired(tt.now) {
					t.Errorf("Verify() should invalidate session by expiring it")
				}
			}
			if tt.name == "error - invalid code" {
				if tt.otp.FailedAttempts() != 1 {
					t.Errorf("Verify() should increment FailedAttempts to 1, got %d", tt.otp.FailedAttempts())
				}
			}
			if tt.name == "error - invalid code blocks after max attempts" {
				if tt.otp.BlockedUntil() == nil {
					t.Error("Verify() should block session after max failed attempts")
				}
				if tt.otp.FailedAttempts() != 3 {
					t.Errorf("Verify() should increment FailedAttempts to 3, got %d", tt.otp.FailedAttempts())
				}
			}
		})
	}
}
