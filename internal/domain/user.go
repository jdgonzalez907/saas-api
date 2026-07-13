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
	id                  int
	personalInformation PersonalInformation
	phone               Phone
	email               *Email
	createdAt           time.Time
	updatedAt           time.Time
}

type UserParams struct {
	ID                  int
	PersonalInformation PersonalInformation
	Phone               Phone
	Email               *Email
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type UserDTO struct {
	ID                     int       `json:"id"`
	PersonalInformationDTO           // Embedded
	Phone                  PhoneDTO  `json:"phone"`
	Email                  *EmailDTO `json:"email"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

func NewUserWithoutId(
	personalInformation PersonalInformation,
	phone Phone,
	email *Email,
) (*User, error) {
	now := time.Now()
	return NewUser(UserParams{
		ID:                  0,
		PersonalInformation: personalInformation,
		Phone:               phone,
		Email:               email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})
}

func NewUser(params UserParams) (*User, error) {
	if params.ID < 0 {
		return nil, ErrInvalidUserID
	}

	return &User{
		id:                  params.ID,
		personalInformation: params.PersonalInformation,
		phone:               params.Phone,
		email:               params.Email,
		createdAt:           params.CreatedAt,
		updatedAt:           params.UpdatedAt,
	}, nil
}

func (u *User) ID() int {
	return u.id
}

func (u *User) Identification() Identification {
	return u.personalInformation.identification
}

func (u *User) FirstName() string {
	return u.personalInformation.firstName
}

func (u *User) LastName() string {
	return u.personalInformation.lastName
}

func (u *User) Phone() Phone {
	return u.phone
}

func (u *User) Email() *Email {
	return u.email
}

func (u *User) Address() *Address {
	return u.personalInformation.address
}

func (u *User) BirthDate() *BirthDate {
	return u.personalInformation.birthDate
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) PersonalInformation() PersonalInformation {
	return u.personalInformation
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

	return &UserDTO{
		ID:                     u.id,
		PersonalInformationDTO: u.personalInformation.ToDTO(),
		Phone:                  u.phone.ToDTO(),
		Email:                  emailDTO,
		CreatedAt:              u.createdAt,
		UpdatedAt:              u.updatedAt,
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
		b, err := NewBirthDate(string(*dto.BirthDate))
		if err != nil {
			return nil, err
		}
		birthDate = &b
	}

	personalInfo, err := NewPersonalInformation(
		identification,
		dto.FirstName,
		dto.LastName,
		address,
		birthDate,
	)
	if err != nil {
		return nil, err
	}

	phone, err := NewPhone(dto.Phone.CountryCode, dto.Phone.Number)
	if err != nil {
		return nil, err
	}

	var email *Email
	if dto.Email != nil {
		e, err := NewEmail(string(*dto.Email))
		if err != nil {
			return nil, err
		}
		email = &e
	}

	return NewUser(UserParams{
		ID:                  dto.ID,
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               email,
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
	})
}

func (u *User) WithPersonalInformation(info PersonalInformation) *User {
	return &User{
		id:                  u.id,
		personalInformation: info,
		phone:               u.phone,
		email:               u.email,
		createdAt:           u.createdAt,
		updatedAt:           time.Now(),
	}
}

func (u *User) WithPhone(phone Phone) *User {
	return &User{
		id:                  u.id,
		personalInformation: u.personalInformation,
		phone:               phone,
		email:               u.email,
		createdAt:           u.createdAt,
		updatedAt:           time.Now(),
	}
}

func (u *User) WithEmail(email *Email) *User {
	return &User{
		id:                  u.id,
		personalInformation: u.personalInformation,
		phone:               u.phone,
		email:               email,
		createdAt:           u.createdAt,
		updatedAt:           time.Now(),
	}
}
