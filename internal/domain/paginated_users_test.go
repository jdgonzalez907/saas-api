package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

func TestPaginatedUsersVO(t *testing.T) {
	now := time.Now()
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(now.AddDate(-18, 0, -1))

	firstPersonalInfo, _ := domain.NewPersonalInformation(identification, "John", "Doe", &address, &birthDate)
	firstUser, _ := domain.NewUser(domain.UserParams{
		ID:                  1,
		PersonalInformation: firstPersonalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	secondPersonalInfo, _ := domain.NewPersonalInformation(identification, "Jane", "Smith", &address, &birthDate)
	secondUser, _ := domain.NewUser(domain.UserParams{
		ID:                  2,
		PersonalInformation: secondPersonalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
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
			users:              []*domain.User{firstUser, secondUser},
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
				actual := dto.Users[0]
				expected := *firstUser.ToDTO()

				if actual.ID != expected.ID {
					t.Errorf("expected ID %d, got %d", expected.ID, actual.ID)
				}
				if actual.FirstName != expected.FirstName {
					t.Errorf("expected FirstName %s, got %s", expected.FirstName, actual.FirstName)
				}
				if actual.LastName != expected.LastName {
					t.Errorf("expected LastName %s, got %s", expected.LastName, actual.LastName)
				}
				if actual.Identification != expected.Identification {
					t.Errorf("expected Identification")
				}
				if actual.Phone != expected.Phone {
					t.Errorf("expected Phone")
				}
				if (actual.Email == nil) != (expected.Email == nil) {
					t.Errorf("expected Email nil mismatch")
				} else if actual.Email != nil {
					if *actual.Email != *expected.Email {
						t.Errorf("expected Email value mismatch")
					}
				}
				if (actual.Address == nil) != (expected.Address == nil) {
					t.Errorf("expected Address nil mismatch")
				} else if actual.Address != nil {
					if *actual.Address != *expected.Address {
						t.Errorf("expected Address value mismatch")
					}
				}
				if (actual.BirthDate == nil) != (expected.BirthDate == nil) {
					t.Errorf("expected BirthDate nil mismatch")
				} else if actual.BirthDate != nil {
					if !actual.BirthDate.Value.Equal(expected.BirthDate.Value) {
						t.Errorf("expected BirthDate value mismatch")
					}
				}
			}
		})
	}
}
