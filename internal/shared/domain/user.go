package domain

type User struct {
	id               int64
	fullName         string
	phoneCountryCode string
	phoneNumber      string
	email            *string
}

type UserDTO struct {
	ID               int64   `json:"id"`
	FullName         string  `json:"full_name"`
	PhoneCountryCode string  `json:"phone_country_code"`
	PhoneNumber      string  `json:"phone_number"`
	Email            *string `json:"email"`
}

func NewUser(id int64, fullName string, phoneCountryCode string, phoneNumber string, email *string) *User {
	return &User{
		id:               id,
		fullName:         fullName,
		phoneCountryCode: phoneCountryCode,
		phoneNumber:      phoneNumber,
		email:            email,
	}
}

func (u *User) ID() int64 {
	return u.id
}

func (u *User) FullName() string {
	return u.fullName
}

func (u *User) PhoneCountryCode() string {
	return u.phoneCountryCode
}

func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

func (u *User) Email() *string {
	return u.email
}

func (u *User) ToDTO() UserDTO {
	return UserDTO{
		ID:               u.id,
		FullName:         u.fullName,
		PhoneCountryCode: u.phoneCountryCode,
		PhoneNumber:      u.phoneNumber,
		Email:            u.email,
	}
}
