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

func TestUpdateUserEmailUseCase(t *testing.T) {
	userID := 1
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-18, 0, -1).Format("2006-01-02"))
	now := time.Now()

	personalInfo, _ := domain.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)

	existingUser, _ := domain.NewUser(domain.UserParams{
		ID:                  userID,
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	otherEmail, _ := domain.NewEmail("other.email@example.com")
	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		inputID          int
		inputEmail       *domain.Email
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName:   "success - no changes to email",
			inputID:    userID,
			inputEmail: &email,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName:   "success - with email change",
			inputID:    userID,
			inputEmail: &otherEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByEmail", mock.Anything, otherEmail).Return(nil, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName:   "success - nil email",
			inputID:    userID,
			inputEmail: nil,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName:   "fail - email already exists",
			inputID:    userID,
			inputEmail: &otherEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByEmail", mock.Anything, otherEmail).Return(existingUser, nil)
			},
			expectedError: domain.ErrUserEmailAlreadyExists,
		},
		{
			testName:   "fail - user not found",
			inputID:    userID,
			inputEmail: &email,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName:   "fail - infra error on FindById",
			inputID:    userID,
			inputEmail: &email,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUserEmail,
		},
		{
			testName:   "fail - infra error on FindByEmail",
			inputID:    userID,
			inputEmail: &otherEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByEmail", mock.Anything, otherEmail).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUserEmail,
		},
		{
			testName:   "fail - repo update error",
			inputID:    userID,
			inputEmail: &email,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(dbErr)
			},
			expectedError: domain.ErrUpdatingUserEmail,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			useCase := application.NewUpdateUserEmailUseCase(mockUserRepository)
			err := useCase.Execute(context.Background(), tc.inputID, tc.inputEmail)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if tc.expectedError == domain.ErrUpdatingUserEmail && errors.Unwrap(err) != nil {
						// Success: wrapped infra error
					} else {
						t.Errorf("expected error: %v, got %v", tc.expectedError, err)
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
