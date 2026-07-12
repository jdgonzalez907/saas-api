package application_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
	domainMocks "jdgonzalez907/users-api/mocks/domain"

	"github.com/stretchr/testify/mock"
)

func TestFindUsersPaginatedUseCase(t *testing.T) {
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

	dbErr := errors.New("db query error")
	cursor := 1

	testCases := []struct {
		testName         string
		pagination       domain.Pagination
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedResult   domain.PaginatedUsers
		expectedError    error
	}{
		{
			testName:   "success - page with users",
			pagination: domain.NewPagination(&cursor, 10),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindAll", mock.MatchedBy(func(p domain.Pagination) bool {
					return p.Limit() == 10 && p.LastID() != nil && *p.LastID() == cursor
				})).Return([]*domain.User{user1, user2}, nil)
			},
			expectedResult: domain.NewPaginatedUsers(
				[]*domain.User{user1, user2},
				func() *int { i := 2; return &i }(),
			),
			expectedError: nil,
		},
		{
			testName:   "success - empty page",
			pagination: domain.NewPagination(nil, 25),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindAll", mock.MatchedBy(func(p domain.Pagination) bool {
					return p.Limit() == 25 && p.LastID() == nil
				})).Return([]*domain.User{}, nil)
			},
			expectedResult: domain.NewPaginatedUsers(
				[]*domain.User{},
				nil,
			),
			expectedError: nil,
		},
		{
			testName:   "fail - repository error",
			pagination: domain.NewPagination(nil, 50),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindAll", mock.MatchedBy(func(p domain.Pagination) bool {
					return p.Limit() == 50 && p.LastID() == nil
				})).Return(nil, dbErr)
			},
			expectedResult: domain.PaginatedUsers{},
			expectedError:  domain.ErrFindingUsers,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			useCase := application.NewFindUsersPaginatedUseCase(mockUserRepository)
			result, err := useCase.Execute(tc.pagination)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if tc.expectedError == domain.ErrFindingUsers && errors.Unwrap(err) != nil {
						// Success: wrapped infra error
					} else {
						t.Errorf("expected error: %v, got %v", tc.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !reflect.DeepEqual(result.Users(), tc.expectedResult.Users()) {
					t.Errorf("expected users: %+v, got %+v", tc.expectedResult.Users(), result.Users())
				}
				if (tc.expectedResult.NextCursor() == nil && result.NextCursor() != nil) ||
					(tc.expectedResult.NextCursor() != nil && result.NextCursor() == nil) ||
					(tc.expectedResult.NextCursor() != nil && result.NextCursor() != nil && *tc.expectedResult.NextCursor() != *result.NextCursor()) {
					t.Errorf("expected next cursor: %v, got %v", tc.expectedResult.NextCursor(), result.NextCursor())
				}
			}
		})
	}
}
