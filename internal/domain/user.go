package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidFirstName = errors.New("invalid first name")
	ErrInvalidLastName  = errors.New("invalid last name")

	ErrUserIDAlreadyExists    = errors.New("user with id already exists")
	ErrUserPhoneAlreadyExists = errors.New("user with phone already exists")
	ErrUserEmailAlreadyExists = errors.New("user with email already exists")
	ErrUserNotFound           = errors.New("user not found")
	ErrCreatingUser           = errors.New("error creating user")
	ErrUpdatingUser           = errors.New("error updating user")
	ErrDeletingUser           = errors.New("error deleting user")
)

type User struct {
	ID             int            `json:"id"`
	Identification Identification `json:"identification"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Phone          Phone          `json:"phone"`
	Email          *Email         `json:"email"`
	Address        *Address       `json:"address"`
	BirthDate      *BirthDate     `json:"birth_date"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func NewUserWithoutId(
	identification Identification,
	firstName, lastName string,
	phone Phone,
	email *Email,
	address *Address,
	birthDate *BirthDate,
) (*User, error) {
	now := time.Now()
	return NewUser(
		0,
		identification,
		firstName,
		lastName,
		phone,
		email,
		address,
		birthDate,
		now,
		now,
	)
}

func NewUser(
	id int,
	identification Identification,
	firstName, lastName string,
	phone Phone,
	email *Email,
	address *Address,
	birthDate *BirthDate,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	if id < 0 {
		return nil, ErrInvalidUserID
	}

	if firstName == "" {
		return nil, ErrInvalidFirstName
	}

	if lastName == "" {
		return nil, ErrInvalidLastName
	}

	return &User{
		ID:             id,
		Identification: identification,
		FirstName:      firstName,
		LastName:       lastName,
		Phone:          phone,
		Email:          email,
		Address:        address,
		BirthDate:      birthDate,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}
