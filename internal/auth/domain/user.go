package domain

import "errors"

var (
	ErrUserInvalidID = errors.New("user ID must be greater than 0")
	ErrUserEmptyName = errors.New("user name cannot be empty")
	ErrUserNotFound  = errors.New("user not found")
)

type User struct {
	id       int64
	fullname string
}

type UserDTO struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
}

func NewUser(id int64, fullname string) (*User, error) {
	if id <= 0 {
		return nil, ErrUserInvalidID
	}
	if fullname == "" {
		return nil, ErrUserEmptyName
	}
	return &User{id: id, fullname: fullname}, nil
}

func (u *User) ID() int64        { return u.id }
func (u *User) Fullname() string { return u.fullname }

func (u *User) Equals(other *User) bool {
	if other == nil {
		return false
	}
	return u.id == other.id
}

func (u *User) ToDTO() UserDTO {
	return UserDTO{ID: u.id, Fullname: u.fullname}
}
