package domain

import (
	"errors"

	"github.com/google/uuid"
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
)

type User struct {
	ID             string         `json:"id"`
	Identification Identification `json:"identification"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Phone          Phone          `json:"phone"`
	Email          *Email         `json:"email"`
	Address        *Address       `json:"address"`
	BirthDate      *BirthDate     `json:"birth_date"`
}

func NewUserWithoutId(
	identification Identification,
	firstName, lastName string,
	phone Phone,
	email *Email,
	address *Address,
	birthDate *BirthDate,
) (*User, error) {
	return NewUser(
		uuid.NewString(),
		identification,
		firstName,
		lastName,
		phone,
		email,
		address,
		birthDate,
	)
}

func NewUser(
	id string,
	identification Identification,
	firstName, lastName string,
	phone Phone,
	email *Email,
	address *Address,
	birthDate *BirthDate,
) (*User, error) {
	if id == "" {
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
	}, nil
}
