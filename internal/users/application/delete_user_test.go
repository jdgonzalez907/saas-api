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

func TestDeleteUserUseCase(t *testing.T) {
	userID := int64(1)
	identification, _ := domain.NewIdentification(domain.IDTypeCC, "1111")
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

	existingUser, _ := domain.NewUser(
		userID,
		personalInfo,
		phone,
		&email,
		now,
		now,
	)

	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		inputID          int64
		executeUserID    int64
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName:      "success - delete user",
			inputID:       userID,
			executeUserID: int64(1),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByID", mock.Anything, userID, mock.Anything).Return(existingUser, nil)
				m.On("Delete", mock.Anything, userID, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName:      "fail - invalid user id",
			inputID:       int64(0),
			executeUserID: int64(1),
			mockExpectations: func(_ *domainMocks.MockUserRepository) {
				// No expectations
			},
			expectedError: domain.ErrInvalidUserID,
		},
		{
			testName:      "fail - user not found",
			inputID:       userID,
			executeUserID: int64(1),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByID", mock.Anything, userID, mock.Anything).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName:      "fail - infra error on FindByID",
			inputID:       userID,
			executeUserID: int64(1),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByID", mock.Anything, userID, mock.Anything).Return(nil, dbErr)
			},
			expectedError: domain.ErrDeletingUser,
		},
		{
			testName:      "fail - infra error on Delete",
			inputID:       userID,
			executeUserID: int64(1),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByID", mock.Anything, userID, mock.Anything).Return(existingUser, nil)
				m.On("Delete", mock.Anything, userID, mock.Anything).Return(dbErr)
			},
			expectedError: domain.ErrDeletingUser,
		},
		{
			testName:      "fail - user ownership mismatch",
			inputID:       userID,
			executeUserID: int64(999),
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByID", mock.Anything, userID, mock.Anything).Return(existingUser, nil)
			},
			expectedError: domain.ErrUserOwnershipMismatch,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			deleteUserUseCase := application.NewDeleteUserUseCase(mockUserRepository)
			err := deleteUserUseCase.Execute(context.Background(), tc.inputID, tc.executeUserID)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrDeletingUser && errors.Unwrap(err) != nil) {
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
