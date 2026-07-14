package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/domain"
	domainMocks "jdgonzalez907/saas-api/mocks/domain"
)

func mustNewPagination(lastID *int64, limitVal *int32) domain.Pagination {
	p, err := domain.NewPagination(lastID, limitVal)
	if err != nil {
		panic(err)
	}
	return p
}

func TestFindUsersPaginatedUseCase(t *testing.T) {
	now := time.Now()
	identification, _ := domain.NewIdentification(domain.IDTypeCC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(now.AddDate(-18, 0, -1).Format("2006-01-02"))

	firstPersonalInfo, _ := domain.NewPersonalInformation(identification, "John", "Doe", &address, &birthDate)
	firstUser, _ := domain.NewUser(domain.UserParams{
		ID:                  int64(1),
		PersonalInformation: firstPersonalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	secondPersonalInfo, _ := domain.NewPersonalInformation(identification, "Jane", "Smith", &address, &birthDate)
	secondUser, _ := domain.NewUser(domain.UserParams{
		ID:                  int64(2),
		PersonalInformation: secondPersonalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	dbErr := errors.New("db query error")
	cursor := int64(1)
	limit10 := int32(10)
	limit25 := int32(25)
	limit50 := int32(50)

	testCases := []struct {
		testName         string
		pagination       domain.Pagination
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedResult   domain.PaginatedUsers
		expectedError    error
	}{
		{
			testName:   "success - page with users",
			pagination: mustNewPagination(&cursor, &limit10),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindAll", mock.Anything, mock.MatchedBy(func(p domain.Pagination) bool {
					return p.Limit() == 10 && p.LastID() != nil && *p.LastID() == cursor
				})).Return([]*domain.User{firstUser, secondUser}, nil)
			},
			expectedResult: domain.NewPaginatedUsers(
				[]*domain.User{firstUser, secondUser},
				func() *int64 { i := int64(2); return &i }(),
			),
			expectedError: nil,
		},
		{
			testName:   "success - empty page",
			pagination: mustNewPagination(nil, &limit25),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindAll", mock.Anything, mock.MatchedBy(func(p domain.Pagination) bool {
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
			pagination: mustNewPagination(nil, &limit50),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindAll", mock.Anything, mock.MatchedBy(func(p domain.Pagination) bool {
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
			result, err := useCase.Execute(context.Background(), tc.pagination)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrFindingUsers && errors.Unwrap(err) != nil) {
						t.Errorf("expected error: %v, got %v", tc.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(result.Users()) != len(tc.expectedResult.Users()) {
					t.Errorf("expected users count: %d, got %d", len(tc.expectedResult.Users()), len(result.Users()))
				} else {
					for idx, expectedUser := range tc.expectedResult.Users() {
						actualUser := result.Users()[idx]
						if actualUser.ID() != expectedUser.ID() {
							t.Errorf("mismatch user ID at index %d: expected %d, got %d", idx, expectedUser.ID(), actualUser.ID())
						}
					}
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
