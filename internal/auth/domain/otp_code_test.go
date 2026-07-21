package domain

import "testing"

func TestNewOTPCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    OTPCode
		wantErr error
	}{
		{
			name:    "success - valid 6 digits",
			code:    "123456",
			want:    "123456",
			wantErr: nil,
		},
		{
			name:    "success - all zeros",
			code:    "000000",
			want:    "000000",
			wantErr: nil,
		},
		{
			name:    "success - max digits",
			code:    "999999",
			want:    "999999",
			wantErr: nil,
		},
		{
			name:    "error - empty code",
			code:    "",
			want:    "",
			wantErr: ErrOTPCodeEmpty,
		},
		{
			name:    "error - too short",
			code:    "12345",
			want:    "",
			wantErr: ErrOTPCodeInvalid,
		},
		{
			name:    "error - too long",
			code:    "1234567",
			want:    "",
			wantErr: ErrOTPCodeInvalid,
		},
		{
			name:    "error - contains letters",
			code:    "12345a",
			want:    "",
			wantErr: ErrOTPCodeInvalid,
		},
		{
			name:    "error - contains special chars",
			code:    "123-456",
			want:    "",
			wantErr: ErrOTPCodeInvalid,
		},
		{
			name:    "error - contains spaces",
			code:    "123 456",
			want:    "",
			wantErr: ErrOTPCodeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOTPCode(tt.code)
			if err != tt.wantErr {
				t.Errorf("NewOTPCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewOTPCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOTPCodeRandom(t *testing.T) {
	code, err := NewOTPCodeRandom()
	if err != nil {
		t.Fatalf("NewOTPCodeRandom() error = %v", err)
	}
	if len(code.Value()) != OTPCodeLength {
		t.Errorf("NewOTPCodeRandom() length = %v, want %v", len(code.Value()), OTPCodeLength)
	}
	for i := 0; i < len(code.Value()); i++ {
		c := code.Value()[i]
		if c < '0' || c > '9' {
			t.Errorf("NewOTPCodeRandom() contains non-digit: %c", c)
		}
	}
}

func TestNewOTPCodeFromString(t *testing.T) {
	v := NewOTPCodeFromString("123456")
	if v != "123456" {
		t.Errorf("NewOTPCodeFromString() = %v, want %v", v, "123456")
	}
}

func TestOTPCode_Value(t *testing.T) {
	v, _ := NewOTPCode("123456")
	if got := v.Value(); got != "123456" {
		t.Errorf("Value() = %v, want %v", got, "123456")
	}
}

func TestOTPCode_Equals(t *testing.T) {
	tests := []struct {
		name string
		v1   OTPCode
		v2   OTPCode
		want bool
	}{
		{
			name: "equal - same value",
			v1:   OTPCode("123456"),
			v2:   OTPCode("123456"),
			want: true,
		},
		{
			name: "not equal - different value",
			v1:   OTPCode("123456"),
			v2:   OTPCode("654321"),
			want: false,
		},
		{
			name: "not equal - one empty",
			v1:   OTPCode("123456"),
			v2:   OTPCode(""),
			want: false,
		},
		{
			name: "equal - both empty",
			v1:   OTPCode(""),
			v2:   OTPCode(""),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Equals(tt.v2); got != tt.want {
				t.Errorf("OTPCode.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOTPCode_ToDTO(t *testing.T) {
	tests := []struct {
		name string
		code OTPCode
		want string
	}{
		{
			name: "success",
			code: OTPCode("123456"),
			want: "123456",
		},
		{
			name: "empty code",
			code: OTPCode(""),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.ToDTO(); got != tt.want {
				t.Errorf("ToDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}
