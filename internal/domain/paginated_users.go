package domain

type PaginatedUsers struct {
	users      []*User
	nextCursor *int
}

func NewPaginatedUsers(users []*User, nextCursor *int) PaginatedUsers {
	return PaginatedUsers{
		users:      users,
		nextCursor: nextCursor,
	}
}

func (p PaginatedUsers) Users() []*User {
	return p.users
}

func (p PaginatedUsers) NextCursor() *int {
	return p.nextCursor
}

type PaginatedUsersDTO struct {
	Users      []UserDTO `json:"users"`
	NextCursor *int      `json:"next_cursor"`
}

func (p PaginatedUsers) ToDTO() *PaginatedUsersDTO {
	userDTOs := make([]UserDTO, len(p.users))
	for i, user := range p.users {
		userDTOs[i] = *user.ToDTO()
	}
	return &PaginatedUsersDTO{
		Users:      userDTOs,
		NextCursor: p.nextCursor,
	}
}
