package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
	domainMocks "jdgonzalez907/users-api/mocks/domain"

	"github.com/stretchr/testify/mock"
)

func TestFindUserByIdUseCase(t *testing.T) {
	now := time.Now()
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(now.AddDate(-18, 0, -1).Format("2006-01-02"))

	personalInfo, _ := domain.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)

	user, _ := domain.NewUser(domain.UserParams{
		ID:                  1,
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	dbErr := errors.New("db query error")

	testCases := []struct {
		testName         string
		inputID          int
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedResult   *domain.User
		expectedError    error
	}{
		{
			testName: "success - user found",
			inputID:  1,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, 1).Return(user, nil)
			},
			expectedResult: user,
			expectedError:  nil,
		},
		{
			testName: "fail - invalid user id (<= 0)",
			inputID:  0,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				// No expectations
			},
			expectedResult: nil,
			expectedError:  domain.ErrInvalidUserID,
		},
		{
			testName: "fail - user not found",
			inputID:  99,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, 99).Return(nil, nil)
			},
			expectedResult: nil,
			expectedError:  domain.ErrUserNotFound,
		},
		{
			testName: "fail - repository error",
			inputID:  1,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, 1).Return(nil, dbErr)
			},
			expectedResult: nil,
			expectedError:  domain.ErrFindingUserByID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			useCase := application.NewFindUserByIDUseCase(mockUserRepository)
			result, err := useCase.Execute(context.Background(), tc.inputID)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if tc.expectedError == domain.ErrFindingUserByID && errors.Unwrap(err) != nil {
						// Success: wrapped infra error
					} else {
						t.Errorf("expected error: %v, got %v", tc.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result != tc.expectedResult {
					t.Errorf("expected result: %+v, got %+v", tc.expectedResult, result)
				}
			}
		})
	}
}
