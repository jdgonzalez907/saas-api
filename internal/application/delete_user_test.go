package application_test

import (
	"errors"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
	domainMocks "jdgonzalez907/users-api/mocks/domain"
)

func TestDeleteUserUseCase(t *testing.T) {
	userID := 1
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-18, 0, -1))
	now := time.Now()

	existingUser, _ := domain.NewUser(domain.UserParams{
		ID:             userID,
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

	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		inputID          int
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName: "success - delete user",
			inputID:  userID,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("Delete", userID).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "fail - user not found",
			inputID:  userID,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName: "fail - infra error on FindById",
			inputID:  userID,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(nil, dbErr)
			},
			expectedError: domain.ErrDeletingUser,
		},
		{
			testName: "fail - infra error on Delete",
			inputID:  userID,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("Delete", userID).Return(dbErr)
			},
			expectedError: domain.ErrDeletingUser,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			testCase.mockExpectations(mockUserRepository)

			deleteUserUseCase := application.NewDeleteUserUseCase(mockUserRepository)
			err := deleteUserUseCase.Execute(testCase.inputID)

			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", testCase.expectedError)
				}
				if !errors.Is(err, testCase.expectedError) {
					if testCase.expectedError == domain.ErrDeletingUser && errors.Unwrap(err) != nil {
						// Success: wrapped infra error
					} else {
						t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}
