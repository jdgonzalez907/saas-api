package domain

import (
	"testing"
	"time"
)

func validBirthdate() time.Time {
	return time.Now().AddDate(-25, 0, 0)
}

func underageBirthdate() time.Time {
	return time.Now().AddDate(-10, 0, 0)
}

func TestNewPersonalInformation(t *testing.T) {
	tests := []struct {
		name      string
		dniType   DNIType
		dniNumber string
		firstName string
		lastName  string
		birthdate time.Time
		wantErr   error
	}{
		{
			name:      "success - cc",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   nil,
		},
		{
			name:      "success - ce",
			dniType:   DNITypeCE,
			dniNumber: "AB1234567",
			firstName: "María",
			lastName:  "García",
			birthdate: validBirthdate(),
			wantErr:   nil,
		},
		{
			name:      "success - nit",
			dniType:   DNITypeNIT,
			dniNumber: "900123456",
			firstName: "Carlos",
			lastName:  "López",
			birthdate: validBirthdate(),
			wantErr:   nil,
		},
		{
			name:      "success - pp",
			dniType:   DNITypePP,
			dniNumber: "PA1234567",
			firstName: "Ana María",
			lastName:  "Martínez",
			birthdate: validBirthdate(),
			wantErr:   nil,
		},
		{
			name:      "error - invalid dni type",
			dniType:   "XX",
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidDNIType,
		},
		{
			name:      "error - empty dni number",
			dniType:   DNITypeCC,
			dniNumber: "",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidDNINumber,
		},
		{
			name:      "error - dni number too short",
			dniType:   DNITypeCC,
			dniNumber: "1234567",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidDNINumber,
		},
		{
			name:      "error - dni number too long",
			dniType:   DNITypeCC,
			dniNumber: "123456789012345678901",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidDNINumber,
		},
		{
			name:      "error - empty first name",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidFirstName,
		},
		{
			name:      "error - first name too short",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Ju",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidFirstName,
		},
		{
			name:      "error - first name too long",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuanjuan",
			lastName:  "Pérez",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidFirstName,
		},
		{
			name:      "error - empty last name",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidLastName,
		},
		{
			name:      "error - last name too short",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "Pé",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidLastName,
		},
		{
			name:      "error - last name too long",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "Garcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcíagarcía",
			birthdate: validBirthdate(),
			wantErr:   ErrPersonalInformationInvalidLastName,
		},
		{
			name:      "error - empty birthdate",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: time.Time{},
			wantErr:   ErrPersonalInformationEmptyBirthdate,
		},
		{
			name:      "error - underage",
			dniType:   DNITypeCC,
			dniNumber: "12345678",
			firstName: "Juan",
			lastName:  "Pérez",
			birthdate: underageBirthdate(),
			wantErr:   ErrPersonalInformationUnderage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPersonalInformation(tt.dniType, tt.dniNumber, tt.firstName, tt.lastName, tt.birthdate)
			if err != tt.wantErr {
				t.Errorf("NewPersonalInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.DNIType() != tt.dniType {
					t.Errorf("DNIType() = %v, want %v", got.DNIType(), tt.dniType)
				}
				if got.DNINumber() != tt.dniNumber {
					t.Errorf("DNINumber() = %v, want %v", got.DNINumber(), tt.dniNumber)
				}
				if got.FirstName() != tt.firstName {
					t.Errorf("FirstName() = %v, want %v", got.FirstName(), tt.firstName)
				}
				if got.LastName() != tt.lastName {
					t.Errorf("LastName() = %v, want %v", got.LastName(), tt.lastName)
				}
				if !got.Birthdate().Equal(tt.birthdate) {
					t.Errorf("Birthdate() = %v, want %v", got.Birthdate(), tt.birthdate)
				}
			}
		})
	}
}

func TestPersonalInformation_DNIType(t *testing.T) {
	v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", validBirthdate())
	if got := v.DNIType(); got != DNITypeCC {
		t.Errorf("DNIType() = %v, want %v", got, DNITypeCC)
	}
}

func TestPersonalInformation_DNINumber(t *testing.T) {
	v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", validBirthdate())
	if got := v.DNINumber(); got != "12345678" {
		t.Errorf("DNINumber() = %v, want %v", got, "12345678")
	}
}

func TestPersonalInformation_FirstName(t *testing.T) {
	v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", validBirthdate())
	if got := v.FirstName(); got != "Juan" {
		t.Errorf("FirstName() = %v, want %v", got, "Juan")
	}
}

func TestPersonalInformation_LastName(t *testing.T) {
	v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", validBirthdate())
	if got := v.LastName(); got != "Pérez" {
		t.Errorf("LastName() = %v, want %v", got, "Pérez")
	}
}

func TestPersonalInformation_Birthdate(t *testing.T) {
	bd := validBirthdate()
	v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
	if got := v.Birthdate(); !got.Equal(bd) {
		t.Errorf("Birthdate() = %v, want %v", got, bd)
	}
}

func TestPersonalInformation_Equals(t *testing.T) {
	bd := validBirthdate()
	tests := []struct {
		name string
		v1   PersonalInformation
		v2   PersonalInformation
		want bool
	}{
		{
			name: "equal - same values",
			v1: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			v2: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			want: true,
		},
		{
			name: "not equal - different dni type",
			v1: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			v2: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCE, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			want: false,
		},
		{
			name: "not equal - different dni number",
			v1: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			v2: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "87654321", "Juan", "Pérez", bd)
				return v
			}(),
			want: false,
		},
		{
			name: "not equal - different first name",
			v1: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			v2: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Carlos", "Pérez", bd)
				return v
			}(),
			want: false,
		},
		{
			name: "not equal - different last name",
			v1: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			v2: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "García", bd)
				return v
			}(),
			want: false,
		},
		{
			name: "not equal - different birthdate",
			v1: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)
				return v
			}(),
			v2: func() PersonalInformation {
				v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd.AddDate(0, 0, 1))
				return v
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Equals(tt.v2); got != tt.want {
				t.Errorf("PersonalInformation.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPersonalInformation_ToDTO(t *testing.T) {
	bd := validBirthdate()
	v, _ := NewPersonalInformation(DNITypeCC, "12345678", "Juan", "Pérez", bd)

	got := v.ToDTO()
	want := PersonalInformationDTO{
		DNIType:   DNITypeCC,
		DNINumber: "12345678",
		FirstName: "Juan",
		LastName:  "Pérez",
		Birthdate: bd.Format(time.RFC3339),
	}

	if got.DNIType != want.DNIType {
		t.Errorf("ToDTO().DNIType = %v, want %v", got.DNIType, want.DNIType)
	}
	if got.DNINumber != want.DNINumber {
		t.Errorf("ToDTO().DNINumber = %v, want %v", got.DNINumber, want.DNINumber)
	}
	if got.FirstName != want.FirstName {
		t.Errorf("ToDTO().FirstName = %v, want %v", got.FirstName, want.FirstName)
	}
	if got.LastName != want.LastName {
		t.Errorf("ToDTO().LastName = %v, want %v", got.LastName, want.LastName)
	}
	if got.Birthdate != want.Birthdate {
		t.Errorf("ToDTO().Birthdate = %v, want %v", got.Birthdate, want.Birthdate)
	}
}
