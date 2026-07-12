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
	ID             int               `json:"id"`
	Identification IdentificationDTO `json:"identification"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	Phone          PhoneDTO          `json:"phone"`
	Email          *EmailDTO         `json:"email"`
	Address        *AddressDTO       `json:"address"`
	BirthDate      *BirthDateDTO     `json:"birth_date"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
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

func (u *User) ID() int {
	return u.id
}

func (u *User) Identification() Identification {
	return u.identification
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Phone() Phone {
	return u.phone
}

func (u *User) Email() *Email {
	return u.email
}

func (u *User) Address() *Address {
	return u.address
}

func (u *User) BirthDate() *BirthDate {
	return u.birthDate
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) ToDTO() *UserDTO {
	if u == nil {
		return nil
	}

	var emailDTO *EmailDTO
	if u.email != nil {
		dto := u.email.ToDTO()
		emailDTO = &dto
	}

	var addressDTO *AddressDTO
	if u.address != nil {
		dto := u.address.ToDTO()
		addressDTO = &dto
	}

	var birthDateDTO *BirthDateDTO
	if u.birthDate != nil {
		dto := u.birthDate.ToDTO()
		birthDateDTO = &dto
	}

	return &UserDTO{
		ID:             u.id,
		Identification: u.identification.ToDTO(),
		FirstName:      u.firstName,
		LastName:       u.lastName,
		Phone:          u.phone.ToDTO(),
		Email:          emailDTO,
		Address:        addressDTO,
		BirthDate:      birthDateDTO,
		CreatedAt:      u.createdAt,
		UpdatedAt:      u.updatedAt,
	}
}

func UserFromDTO(dto *UserDTO) (*User, error) {
	if dto == nil {
		return nil, nil
	}

	identification, err := NewIdentification(dto.Identification.Type, dto.Identification.Number)
	if err != nil {
		return nil, err
	}

	phone, err := NewPhone(dto.Phone.Value)
	if err != nil {
		return nil, err
	}

	var email *Email
	if dto.Email != nil {
		e, err := NewEmail(dto.Email.Value)
		if err != nil {
			return nil, err
		}
		email = &e
	}

	var address *Address
	if dto.Address != nil {
		a, err := NewAddress(
			dto.Address.Street,
			dto.Address.City,
			dto.Address.State,
			dto.Address.Country,
			dto.Address.PostalCode,
			dto.Address.Description,
		)
		if err != nil {
			return nil, err
		}
		address = &a
	}

	var birthDate *BirthDate
	if dto.BirthDate != nil {
		b, err := NewBirthDate(dto.BirthDate.Value)
		if err != nil {
			return nil, err
		}
		birthDate = &b
	}

	return NewUser(UserParams{
		ID:             dto.ID,
		Identification: identification,
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		Phone:          phone,
		Email:          email,
		Address:        address,
		BirthDate:      birthDate,
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
