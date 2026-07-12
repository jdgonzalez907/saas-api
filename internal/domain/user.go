package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidFirstName = errors.New("invalid first name")
	ErrInvalidLastName  = errors.New("invalid last name")

	ErrUserIDAlreadyExists             = errors.New("user with id already exists")
	ErrUserPhoneAlreadyExists          = errors.New("user with phone already exists")
	ErrUserEmailAlreadyExists          = errors.New("user with email already exists")
	ErrUserNotFound                    = errors.New("user not found")
	ErrCreatingUser                    = errors.New("error creating user")
	ErrUpdatingUserPersonalInformation = errors.New("error updating user personal information")
	ErrUpdatingUserPhone               = errors.New("error updating user phone")
	ErrUpdatingUserEmail               = errors.New("error updating user email")
	ErrDeletingUser                    = errors.New("error deleting user")
	ErrFindingUsers                    = errors.New("error finding users")
	ErrFindingUserByID                 = errors.New("error finding user by id")
)

type User struct {
	id             int
	identification Identification
	firstName      string
	lastName       string
	phone          Phone
	email          *Email
	address        *Address
	birthDate      *BirthDate
	createdAt      time.Time
	updatedAt      time.Time
}

type UserParams struct {
	ID             int
	Identification Identification
	FirstName      string
	LastName       string
	Phone          Phone
	Email          *Email
	Address        *Address
	BirthDate      *BirthDate
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserDTO struct {
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
	return NewUser(UserParams{
		ID:             0,
		Identification: identification,
		FirstName:      firstName,
		LastName:       lastName,
		Phone:          phone,
		Email:          email,
		Address:        address,
		BirthDate:      birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
}

func NewUser(params UserParams) (*User, error) {
	if params.ID < 0 {
		return nil, ErrInvalidUserID
	}

	if params.FirstName == "" {
		return nil, ErrInvalidFirstName
	}

	if params.LastName == "" {
		return nil, ErrInvalidLastName
	}

	return &User{
		id:             params.ID,
		identification: params.Identification,
		firstName:      params.FirstName,
		lastName:       params.LastName,
		phone:          params.Phone,
		email:          params.Email,
		address:        params.Address,
		birthDate:      params.BirthDate,
		createdAt:      params.CreatedAt,
		updatedAt:      params.UpdatedAt,
	}, nil
}

func (u *User) ToDTO() *UserDTO {
	if u == nil {
		return nil
	}
	return &UserDTO{
		ID:             u.id,
		Identification: u.identification,
		FirstName:      u.firstName,
		LastName:       u.lastName,
		Phone:          u.phone,
		Email:          u.email,
		Address:        u.address,
		BirthDate:      u.birthDate,
		CreatedAt:      u.createdAt,
		UpdatedAt:      u.updatedAt,
	}
}

func UserFromDTO(dto *UserDTO) (*User, error) {
	if dto == nil {
		return nil, nil
	}
	return NewUser(UserParams{
		ID:             dto.ID,
		Identification: dto.Identification,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		Phone:          dto.Phone,
		Email:          dto.Email,
		Address:        dto.Address,
		BirthDate:      dto.BirthDate,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	})
}

func (u *User) WithPersonalInformation(
	identification Identification,
	firstName, lastName string,
	address *Address,
	birthDate *BirthDate,
) (*User, error) {
	if firstName == "" {
		return nil, ErrInvalidFirstName
	}

	if lastName == "" {
		return nil, ErrInvalidLastName
	}

	return &User{
		id:             u.id,
		identification: identification,
		firstName:      firstName,
		lastName:       lastName,
		phone:          u.phone,
		email:          u.email,
		address:        address,
		birthDate:      birthDate,
		createdAt:      u.createdAt,
		updatedAt:      time.Now(),
	}, nil
}

func (u *User) WithPhone(phone Phone) *User {
	return &User{
		id:             u.id,
		identification: u.identification,
		firstName:      u.firstName,
		lastName:       u.lastName,
		phone:          phone,
		email:          u.email,
		address:        u.address,
		birthDate:      u.birthDate,
		createdAt:      u.createdAt,
		updatedAt:      time.Now(),
	}
}

func (u *User) WithEmail(email *Email) *User {
	return &User{
		id:             u.id,
		identification: u.identification,
		firstName:      u.firstName,
		lastName:       u.lastName,
		phone:          u.phone,
		email:          email,
		address:        u.address,
		birthDate:      u.birthDate,
		createdAt:      u.createdAt,
		updatedAt:      time.Now(),
	}
}
