package domain

import "testing"

func TestNewPhone(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		number      string
		wantErr     error
	}{
		{
			name:        "success - valid phone number",
			countryCode: "57",
			number:      "3001234567",
			wantErr:     nil,
		},
		{
			name:        "success - country code 1 digit",
			countryCode: "1",
			number:      "2025551234",
			wantErr:     nil,
		},
		{
			name:        "success - country code 3 digits",
			countryCode: "573",
			number:      "3001234567",
			wantErr:     nil,
		},
		{
			name:        "success - number 15 digits",
			countryCode: "57",
			number:      "300123456789012",
			wantErr:     nil,
		},
		{
			name:        "error - empty country code",
			countryCode: "",
			number:      "3001234567",
			wantErr:     ErrPhoneEmptyCountryCode,
		},
		{
			name:        "error - country code too long",
			countryCode: "1234",
			number:      "3001234567",
			wantErr:     ErrPhoneInvalidCountryCode,
		},
		{
			name:        "error - country code non digits",
			countryCode: "abc",
			number:      "3001234567",
			wantErr:     ErrPhoneNonDigitCountryCode,
		},
		{
			name:        "error - empty number",
			countryCode: "57",
			number:      "",
			wantErr:     ErrPhoneEmptyNumber,
		},
		{
			name:        "error - number too short",
			countryCode: "57",
			number:      "123456789",
			wantErr:     ErrPhoneInvalidNumber,
		},
		{
			name:        "error - number too long",
			countryCode: "57",
			number:      "1234567890123456",
			wantErr:     ErrPhoneInvalidNumber,
		},
		{
			name:        "error - number non digits",
			countryCode: "57",
			number:      "abc1234567890",
			wantErr:     ErrPhoneNonDigitNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPhone(tt.countryCode, tt.number)
			if err != tt.wantErr {
				t.Errorf("NewPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.CountryCode() != tt.countryCode {
					t.Errorf("CountryCode() = %v, want %v", got.CountryCode(), tt.countryCode)
				}
				if got.Number() != tt.number {
					t.Errorf("Number() = %v, want %v", got.Number(), tt.number)
				}
			}
		})
	}
}

func TestPhone_CountryCode(t *testing.T) {
	v, _ := NewPhone("57", "3001234567")
	if got := v.CountryCode(); got != "57" {
		t.Errorf("CountryCode() = %v, want %v", got, "57")
	}
}

func TestPhone_Number(t *testing.T) {
	v, _ := NewPhone("57", "3001234567")
	if got := v.Number(); got != "3001234567" {
		t.Errorf("Number() = %v, want %v", got, "3001234567")
	}
}

func TestPhone_Equals(t *testing.T) {
	tests := []struct {
		name string
		v1   Phone
		v2   Phone
		want bool
	}{
		{
			name: "equal - same values",
			v1: func() Phone {
				v, _ := NewPhone("57", "3001234567")
				return v
			}(),
			v2: func() Phone {
				v, _ := NewPhone("57", "3001234567")
				return v
			}(),
			want: true,
		},
		{
			name: "not equal - different country code",
			v1: func() Phone {
				v, _ := NewPhone("57", "3001234567")
				return v
			}(),
			v2: func() Phone {
				v, _ := NewPhone("1", "3001234567")
				return v
			}(),
			want: false,
		},
		{
			name: "not equal - different number",
			v1: func() Phone {
				v, _ := NewPhone("57", "3001234567")
				return v
			}(),
			v2: func() Phone {
				v, _ := NewPhone("57", "3009876543")
				return v
			}(),
			want: false,
		},
		{
			name: "not equal - both different",
			v1: func() Phone {
				v, _ := NewPhone("57", "3001234567")
				return v
			}(),
			v2: func() Phone {
				v, _ := NewPhone("1", "2025551234")
				return v
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Equals(tt.v2); got != tt.want {
				t.Errorf("Phone.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhone_ToDTO(t *testing.T) {
	v, _ := NewPhone("57", "3001234567")
	got := v.ToDTO()
	want := PhoneDTO{
		CountryCode: "57",
		Number:      "3001234567",
	}

	if got.CountryCode != want.CountryCode {
		t.Errorf("ToDTO().CountryCode = %v, want %v", got.CountryCode, want.CountryCode)
	}
	if got.Number != want.Number {
		t.Errorf("ToDTO().Number = %v, want %v", got.Number, want.Number)
	}
}
