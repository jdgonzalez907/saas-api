package domain_test

import (
	"reflect"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

func TestPaginatedUsersVO(t *testing.T) {
	now := time.Now()
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(now.AddDate(-18, 0, -1))

	user1, _ := domain.NewUser(domain.UserParams{
		ID:             1,
		Identification: identification,
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	user2, _ := domain.NewUser(domain.UserParams{
		ID:             2,
		Identification: identification,
		FirstName:      "Jane",
		LastName:       "Smith",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	nextCursor := 2

	testCases := []struct {
		testName           string
		users              []*domain.User
		nextCursor         *int
		expectedUsersCount int
	}{
		{
			testName:           "success - create paginated users VO and convert to DTO",
			users:              []*domain.User{user1, user2},
			nextCursor:         &nextCursor,
			expectedUsersCount: 2,
		},
		{
			testName:           "success - empty paginated users VO",
			users:              []*domain.User{},
			nextCursor:         nil,
			expectedUsersCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			vo := domain.NewPaginatedUsers(tc.users, tc.nextCursor)

			if len(vo.Users()) != tc.expectedUsersCount {
				t.Errorf("expected count: %d, got: %d", tc.expectedUsersCount, len(vo.Users()))
			}

			if vo.NextCursor() != tc.nextCursor {
				t.Errorf("expected cursor: %v, got: %v", tc.nextCursor, vo.NextCursor())
			}

			dto := vo.ToDTO()
			if len(dto.Users) != tc.expectedUsersCount {
				t.Errorf("expected DTO count: %d, got: %d", tc.expectedUsersCount, len(dto.Users))
			}

			if dto.NextCursor != tc.nextCursor {
				t.Errorf("expected DTO cursor: %v, got: %v", tc.nextCursor, dto.NextCursor)
			}

			if tc.expectedUsersCount > 0 {
				if !reflect.DeepEqual(dto.Users[0], *user1.ToDTO()) {
					t.Errorf("expected DTO user1: %+v, got %+v", *user1.ToDTO(), dto.Users[0])
				}
			}
		})
	}
}
