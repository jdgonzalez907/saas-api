package domain

import "testing"

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    Email
		wantErr error
	}{
		{
			name:    "success - valid email",
			email:   "user@example.com",
			want:    "user@example.com",
			wantErr: nil,
		},
		{
			name:    "success - email with subdomain",
			email:   "user@sub.example.com",
			want:    "user@sub.example.com",
			wantErr: nil,
		},
		{
			name:    "error - empty email",
			email:   "",
			want:    "",
			wantErr: ErrEmailEmpty,
		},
		{
			name:    "error - invalid format",
			email:   "not-an-email",
			want:    "",
			wantErr: ErrEmailInvalid,
		},
		{
			name:    "error - missing @",
			email:   "userexample.com",
			want:    "",
			wantErr: ErrEmailInvalid,
		},
		{
			name:    "error - missing domain",
			email:   "user@",
			want:    "",
			wantErr: ErrEmailInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.email)
			if err != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name string
		v1   Email
		v2   Email
		want bool
	}{
		{
			name: "equal - same email",
			v1:   "user@example.com",
			v2:   "user@example.com",
			want: true,
		},
		{
			name: "equal - same value different case",
			v1:   "User@Example.com",
			v2:   "user@example.com",
			want: false,
		},
		{
			name: "not equal - different emails",
			v1:   "user1@example.com",
			v2:   "user2@example.com",
			want: false,
		},
		{
			name: "not equal - one empty",
			v1:   "user@example.com",
			v2:   "",
			want: false,
		},
		{
			name: "equal - both empty",
			v1:   "",
			v2:   "",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Equals(tt.v2); got != tt.want {
				t.Errorf("Email.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_ToDTO(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{
			name:  "success",
			email: "user@example.com",
			want:  "user@example.com",
		},
		{
			name:  "empty email",
			email: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.email.ToDTO(); got != tt.want {
				t.Errorf("Email.ToDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}
