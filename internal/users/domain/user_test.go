package domain

import (
	"testing"
	"time"
)

func validPersonalInformation(t *testing.T) PersonalInformation {
	t.Helper()
	pi, err := NewPersonalInformation(DNITypeCC, "1234567890", "John", "DoeSmith", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("validPersonalInformation() error = %v", err)
	}
	return pi
}

func validPhone(t *testing.T) Phone {
	t.Helper()
	p, err := NewPhone("57", "3001234567")
	if err != nil {
		t.Fatalf("validPhone() error = %v", err)
	}
	return p
}

func validEmail(t *testing.T) *Email {
	t.Helper()
	e, err := NewEmail("john@example.com")
	if err != nil {
		t.Fatalf("validEmail() error = %v", err)
	}
	return &e
}

func TestNew(t *testing.T) {
	email := validEmail(t)
	pi := validPersonalInformation(t)
	phone := validPhone(t)

	tests := []struct {
		name                string
		email               *Email
		personalInformation PersonalInformation
		phone               Phone
	}{
		{
			name:                "success - with email",
			email:               email,
			personalInformation: pi,
			phone:               phone,
		},
		{
			name:                "success - without email",
			email:               nil,
			personalInformation: pi,
			phone:               phone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := New(tt.email, tt.personalInformation, tt.phone)
			if err != nil {
				t.Errorf("New() error = %v", err)
				return
			}
			if u.ID() != 0 {
				t.Errorf("New().ID() = %v, want %v", u.ID(), 0)
			}
			if u.Email() != tt.email {
				t.Errorf("New().Email() = %v, want %v", u.Email(), tt.email)
			}
			if u.PersonalInformation() != tt.personalInformation {
				t.Errorf("New().PersonalInformation() = %v, want %v", u.PersonalInformation(), tt.personalInformation)
			}
			if u.Phone() != tt.phone {
				t.Errorf("New().Phone() = %v, want %v", u.Phone(), tt.phone)
			}
			if u.CreatedAt().IsZero() {
				t.Errorf("New().CreatedAt() should not be zero")
			}
			if u.UpdatedAt().IsZero() {
				t.Errorf("New().UpdatedAt() should not be zero")
			}
		})
	}
}

func TestNewWithID(t *testing.T) {
	email := validEmail(t)
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	now := time.Now()

	tests := []struct {
		name                string
		id                  int64
		email               *Email
		personalInformation PersonalInformation
		phone               Phone
		createdAt           time.Time
		updatedAt           time.Time
		wantErr             error
	}{
		{
			name:                "success - with email",
			id:                  1,
			email:               email,
			personalInformation: pi,
			phone:               phone,
			createdAt:           now,
			updatedAt:           now,
			wantErr:             nil,
		},
		{
			name:                "success - without email",
			id:                  1,
			email:               nil,
			personalInformation: pi,
			phone:               phone,
			createdAt:           now,
			updatedAt:           now,
			wantErr:             nil,
		},
		{
			name:                "error - zero ID",
			id:                  0,
			email:               email,
			personalInformation: pi,
			phone:               phone,
			createdAt:           now,
			updatedAt:           now,
			wantErr:             ErrUserInvalidID,
		},
		{
			name:                "error - negative ID",
			id:                  -1,
			email:               email,
			personalInformation: pi,
			phone:               phone,
			createdAt:           now,
			updatedAt:           now,
			wantErr:             ErrUserInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewWithID(tt.id, tt.email, tt.personalInformation, tt.phone, tt.createdAt, tt.updatedAt)
			if err != tt.wantErr {
				t.Errorf("NewWithID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if u.ID() != tt.id {
					t.Errorf("NewWithID().ID() = %v, want %v", u.ID(), tt.id)
				}
				if u.Email() != tt.email {
					t.Errorf("NewWithID().Email() = %v, want %v", u.Email(), tt.email)
				}
				if u.PersonalInformation() != tt.personalInformation {
					t.Errorf("NewWithID().PersonalInformation() = %v, want %v", u.PersonalInformation(), tt.personalInformation)
				}
				if u.Phone() != tt.phone {
					t.Errorf("NewWithID().Phone() = %v, want %v", u.Phone(), tt.phone)
				}
				if !u.CreatedAt().Equal(tt.createdAt) {
					t.Errorf("NewWithID().CreatedAt() = %v, want %v", u.CreatedAt(), tt.createdAt)
				}
				if !u.UpdatedAt().Equal(tt.updatedAt) {
					t.Errorf("NewWithID().UpdatedAt() = %v, want %v", u.UpdatedAt(), tt.updatedAt)
				}
			}
		})
	}
}

func TestUser_Equals(t *testing.T) {
	email := validEmail(t)
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	now := time.Now()

	tests := []struct {
		name string
		e1   *User
		e2   *User
		want bool
	}{
		{
			name: "equal - same ID",
			e1:   mustNewUserWithID(t, 1, email, pi, phone, now, now),
			e2:   mustNewUserWithID(t, 1, email, pi, phone, now, now),
			want: true,
		},
		{
			name: "not equal - different ID",
			e1:   mustNewUserWithID(t, 1, email, pi, phone, now, now),
			e2:   mustNewUserWithID(t, 2, email, pi, phone, now, now),
			want: false,
		},
		{
			name: "not equal - nil other",
			e1:   mustNewUserWithID(t, 1, email, pi, phone, now, now),
			e2:   nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e1.Equals(tt.e2); got != tt.want {
				t.Errorf("User.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ChangeEmail(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	now := time.Now()
	newEmail := Email("new@example.com")
	sameEmail := Email("john@example.com")

	tests := []struct {
		name       string
		entity     *User
		email      Email
		modifiedBy int64
		wantErr    error
	}{
		{
			name:       "success - change to new email",
			entity:     mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now),
			email:      newEmail,
			modifiedBy: 1,
			wantErr:    nil,
		},
		{
			name:       "success - set email when nil",
			entity:     mustNewUserWithID(t, 1, nil, pi, phone, now, now),
			email:      newEmail,
			modifiedBy: 1,
			wantErr:    nil,
		},
		{
			name:       "success - same email early return",
			entity:     mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now),
			email:      sameEmail,
			modifiedBy: 1,
			wantErr:    nil,
		},
		{
			name:       "error - unauthorized modification",
			entity:     mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now),
			email:      newEmail,
			modifiedBy: 2,
			wantErr:    ErrUserUnauthorizedModification,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := tt.entity.UpdatedAt()
			err := tt.entity.ChangeEmail(tt.email, tt.modifiedBy)
			if err != tt.wantErr {
				t.Errorf("User.ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got := tt.entity.Email(); got == nil || !got.Equals(tt.email) {
					t.Errorf("User.Email() = %v, want %v", got, tt.email)
				}
				if tt.entity.UpdatedAt().Before(updatedAt) {
					t.Errorf("User.UpdatedAt() should be updated")
				}
			}
		})
	}
}

func TestUser_UpdatePersonalInformation(t *testing.T) {
	email := validEmail(t)
	phone := validPhone(t)
	now := time.Now()
	newPI, _ := NewPersonalInformation(DNITypeCE, "AB9876543", "Jane", "Smith", time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name                string
		entity              *User
		personalInformation PersonalInformation
		modifiedBy          int64
		wantErr             error
	}{
		{
			name:                "success",
			entity:              mustNewUserWithID(t, 1, email, validPersonalInformation(t), phone, now, now),
			personalInformation: newPI,
			modifiedBy:          1,
			wantErr:             nil,
		},
		{
			name:                "success - same personal information early return",
			entity:              mustNewUserWithID(t, 1, email, validPersonalInformation(t), phone, now, now),
			personalInformation: validPersonalInformation(t),
			modifiedBy:          1,
			wantErr:             nil,
		},
		{
			name:                "error - unauthorized modification",
			entity:              mustNewUserWithID(t, 1, email, validPersonalInformation(t), phone, now, now),
			personalInformation: newPI,
			modifiedBy:          2,
			wantErr:             ErrUserUnauthorizedModification,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := tt.entity.UpdatedAt()
			err := tt.entity.UpdatePersonalInformation(tt.personalInformation, tt.modifiedBy)
			if err != tt.wantErr {
				t.Errorf("User.UpdatePersonalInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got := tt.entity.PersonalInformation(); got != tt.personalInformation {
					t.Errorf("User.PersonalInformation() = %v, want %v", got, tt.personalInformation)
				}
				if tt.entity.UpdatedAt().Before(updatedAt) {
					t.Errorf("User.UpdatedAt() should be updated")
				}
			}
		})
	}
}

func TestUser_ChangePhone(t *testing.T) {
	email := validEmail(t)
	pi := validPersonalInformation(t)
	now := time.Now()
	newPhone, _ := NewPhone("1", "2025551234")

	tests := []struct {
		name       string
		entity     *User
		phone      Phone
		modifiedBy int64
		wantErr    error
	}{
		{
			name:       "success",
			entity:     mustNewUserWithID(t, 1, email, pi, validPhone(t), now, now),
			phone:      newPhone,
			modifiedBy: 1,
			wantErr:    nil,
		},
		{
			name:       "success - same phone early return",
			entity:     mustNewUserWithID(t, 1, email, pi, validPhone(t), now, now),
			phone:      validPhone(t),
			modifiedBy: 1,
			wantErr:    nil,
		},
		{
			name:       "error - unauthorized modification",
			entity:     mustNewUserWithID(t, 1, email, pi, validPhone(t), now, now),
			phone:      newPhone,
			modifiedBy: 2,
			wantErr:    ErrUserUnauthorizedModification,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := tt.entity.UpdatedAt()
			err := tt.entity.ChangePhone(tt.phone, tt.modifiedBy)
			if err != tt.wantErr {
				t.Errorf("User.ChangePhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got := tt.entity.Phone(); got != tt.phone {
					t.Errorf("User.Phone() = %v, want %v", got, tt.phone)
				}
				if tt.entity.UpdatedAt().Before(updatedAt) {
					t.Errorf("User.UpdatedAt() should be updated")
				}
			}
		})
	}
}

func TestUser_ToDTO(t *testing.T) {
	email := validEmail(t)
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	now := time.Now()
	u := mustNewUserWithID(t, 1, email, pi, phone, now, now)

	dto := u.ToDTO()

	if dto.ID != 1 {
		t.Errorf("User.ToDTO().ID = %v, want %v", dto.ID, 1)
	}

	if dto.Email == nil {
		t.Errorf("User.ToDTO().Email should not be nil")
	} else if *dto.Email != email.ToDTO() {
		t.Errorf("User.ToDTO().Email = %v, want %v", *dto.Email, email.ToDTO())
	}

	if dto.PersonalInformation != pi.ToDTO() {
		t.Errorf("User.ToDTO().PersonalInformation = %v, want %v", dto.PersonalInformation, pi.ToDTO())
	}

	if dto.Phone != phone.ToDTO() {
		t.Errorf("User.ToDTO().Phone = %v, want %v", dto.Phone, phone.ToDTO())
	}

	if !dto.CreatedAt.Equal(now) {
		t.Errorf("User.ToDTO().CreatedAt = %v, want %v", dto.CreatedAt, now)
	}

	if !dto.UpdatedAt.Equal(now) {
		t.Errorf("User.ToDTO().UpdatedAt = %v, want %v", dto.UpdatedAt, now)
	}
}

func TestUser_ToDTO_NilEmail(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	now := time.Now()
	u := mustNewUserWithID(t, 1, nil, pi, phone, now, now)

	dto := u.ToDTO()

	if dto.Email != nil {
		t.Errorf("User.ToDTO().Email should be nil, got %v", *dto.Email)
	}
}

func TestUser_AssignID(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	u := mustNewUser(t, validEmail(t), pi, phone)

	u.AssignID(42)

	if got := u.ID(); got != 42 {
		t.Errorf("User.ID() = %v, want %v", got, 42)
	}
}

func mustNewUser(t *testing.T, email *Email, pi PersonalInformation, phone Phone) *User {
	t.Helper()
	u, err := New(email, pi, phone)
	if err != nil {
		t.Fatalf("mustNewUser() error = %v", err)
	}
	return u
}

func mustNewUserWithID(t *testing.T, id int64, email *Email, pi PersonalInformation, phone Phone, createdAt, updatedAt time.Time) *User {
	t.Helper()
	u, err := NewWithID(id, email, pi, phone, createdAt, updatedAt)
	if err != nil {
		t.Fatalf("mustNewUserWithID() error = %v", err)
	}
	return u
}
