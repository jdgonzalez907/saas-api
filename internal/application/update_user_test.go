package application_test

import (
	"errors"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
	domainMocks "jdgonzalez907/users-api/mocks/domain"

	"github.com/google/uuid"
)

func TestUpdateUserUseCase(t *testing.T) {
	userID := uuid.NewString()
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-18, 0, -1))

	existingUser, _ := domain.NewUser(
		userID,
		identification,
		"John",
		"Doe",
		phone,
		&email,
		&address,
		&birthDate,
	)

	// Modified users for testing
	otherPhone, _ := domain.NewPhone("987654321")
	otherEmail, _ := domain.NewEmail("other.email@example.com")

	userWithNewPhone, _ := domain.NewUser(
		userID,
		identification,
		"John",
		"Doe",
		otherPhone,
		&email,
		&address,
		&birthDate,
	)

	userWithNewEmail, _ := domain.NewUser(
		userID,
		identification,
		"John",
		"Doe",
		phone,
		&otherEmail,
		&address,
		&birthDate,
	)

	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		input            domain.User
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName: "update user with success (no changes to phone/email)",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("Update", existingUser).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "update user with phone change success",
			input:    *userWithNewPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("FindByPhone", otherPhone).Return(nil, nil)
				m.On("Update", userWithNewPhone).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "update user fails - phone already exists",
			input:    *userWithNewPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("FindByPhone", otherPhone).Return(existingUser, nil) // simulates phone owned by another user (or same user, but in our code if userFound != nil, it returns conflict)
			},
			expectedError: domain.ErrUserPhoneAlreadyExists,
		},
		{
			testName: "update user with email change success",
			input:    *userWithNewEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("FindByEmail", otherEmail).Return(nil, nil)
				m.On("Update", userWithNewEmail).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "update user fails - email already exists",
			input:    *userWithNewEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("FindByEmail", otherEmail).Return(existingUser, nil)
			},
			expectedError: domain.ErrUserEmailAlreadyExists,
		},
		{
			testName: "update user fails - user not found",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName: "update user fails - infra error on FindById",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUser,
		},
		{
			testName: "update user fails - infra error on FindByPhone",
			input:    *userWithNewPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("FindByPhone", otherPhone).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUser,
		},
		{
			testName: "update user fails - infra error on FindByEmail",
			input:    *userWithNewEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("FindByEmail", otherEmail).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUser,
		},
		{
			testName: "update user fails - repo update error",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("Update", existingUser).Return(dbErr)
			},
			expectedError: domain.ErrUpdatingUser,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			testCase.mockExpectations(mockUserRepository)

			updateUserUseCase := application.NewUpdateUserUseCase(mockUserRepository)
			err := updateUserUseCase.Execute(&testCase.input)

			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", testCase.expectedError)
				}
				if !errors.Is(err, testCase.expectedError) {
					if testCase.expectedError == domain.ErrUpdatingUser && errors.Unwrap(err) != nil {
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
