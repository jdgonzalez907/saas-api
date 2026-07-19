package domain

import (
	"errors"
	"time"
)

var (
	ErrUserInvalidID                = errors.New("user ID must be greater than 0")
	ErrUserUnauthorizedModification = errors.New("user is not authorized to modify this resource")
	ErrUserPhoneAlreadyExists       = errors.New("user with this phone already exists")
	ErrUserEmailAlreadyExists       = errors.New("user with this email already exists")

	ErrUserNotFound = errors.New("user not found")
	ErrFindUserByID = errors.New("cannot find user by ID")

	ErrCreateUser                = errors.New("cannot create user")
	ErrChangePhone               = errors.New("cannot change phone")
	ErrChangeEmail               = errors.New("cannot change email")
	ErrUpdatePersonalInformation = errors.New("cannot update personal information")
)

type User struct {
	id                  int64
	email               *Email
	personalInformation PersonalInformation
	phone               Phone
	createdAt           time.Time
	updatedAt           time.Time
}

type UserDTO struct {
	ID                  int64                  `json:"id"`
	Email               *string                `json:"email"`
	PersonalInformation PersonalInformationDTO `json:"personal_information"`
	Phone               PhoneDTO               `json:"phone"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// New creates a new User for creation (ID auto-generated later)
func New(email *Email, personalInformation PersonalInformation, phone Phone) (*User, error) {
	now := time.Now()
	return &User{
		email:               email,
		personalInformation: personalInformation,
		phone:               phone,
		createdAt:           now,
		updatedAt:           now,
	}, nil
}

// NewWithID creates a User with ID (for loading from DB)
func NewWithID(id int64, email *Email, personalInformation PersonalInformation, phone Phone, createdAt, updatedAt time.Time) (*User, error) {
	if id <= 0 {
		return nil, ErrUserInvalidID
	}

	return &User{
		id:                  id,
		email:               email,
		personalInformation: personalInformation,
		phone:               phone,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
	}, nil
}

func (u *User) ID() int64                                { return u.id }
func (u *User) Email() *Email                            { return u.email }
func (u *User) PersonalInformation() PersonalInformation { return u.personalInformation }
func (u *User) Phone() Phone                             { return u.phone }
func (u *User) CreatedAt() time.Time                     { return u.createdAt }
func (u *User) UpdatedAt() time.Time                     { return u.updatedAt }

func (u *User) AssignID(id int64) {
	u.id = id
}

func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id == other.id
}

func (u *User) ChangeEmail(email Email, modifiedBy int64) error {
	if modifiedBy != u.id {
		return ErrUserUnauthorizedModification
	}

	if u.email != nil && u.email.Equals(email) {
		return nil
	}

	u.email = &email
	u.updatedAt = time.Now()
	return nil
}

func (u *User) UpdatePersonalInformation(pi PersonalInformation, modifiedBy int64) error {
	if modifiedBy != u.id {
		return ErrUserUnauthorizedModification
	}

	if u.personalInformation.Equals(pi) {
		return nil
	}

	u.personalInformation = pi
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangePhone(phone Phone, modifiedBy int64) error {
	if modifiedBy != u.id {
		return ErrUserUnauthorizedModification
	}

	if u.phone.Equals(phone) {
		return nil
	}

	u.phone = phone
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ToDTO() UserDTO {
	var emailStr *string
	if u.email != nil {
		s := u.email.ToDTO()
		emailStr = &s
	}

	return UserDTO{
		ID:                  u.id,
		Email:               emailStr,
		PersonalInformation: u.personalInformation.ToDTO(),
		Phone:               u.phone.ToDTO(),
		CreatedAt:           u.createdAt,
		UpdatedAt:           u.updatedAt,
	}
}
