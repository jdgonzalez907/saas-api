package domain

type UserRepository interface {
	FindById(id string) (*User, error)
	FindByPhone(phone Phone) (*User, error)
	FindByEmail(email Email) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}
