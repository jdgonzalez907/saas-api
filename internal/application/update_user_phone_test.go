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

func TestUpdateUserPhoneUseCase(t *testing.T) {
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

	otherPhone, _ := domain.NewPhone("57", "987654321")
	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		inputID          int
		inputPhone       domain.Phone
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName:   "success - no changes to phone",
			inputID:    userID,
			inputPhone: phone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
			},
			expectedError: nil,
		},
		{
			testName:   "success - with phone change",
			inputID:    userID,
			inputPhone: otherPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByPhone", mock.Anything, otherPhone).Return(nil, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName:   "fail - phone already exists",
			inputID:    userID,
			inputPhone: otherPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByPhone", mock.Anything, otherPhone).Return(existingUser, nil)
			},
			expectedError: domain.ErrUserPhoneAlreadyExists,
		},
		{
			testName:   "fail - user not found",
			inputID:    userID,
			inputPhone: phone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName:   "fail - infra error on FindById",
			inputID:    userID,
			inputPhone: phone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUserPhone,
		},
		{
			testName:   "fail - infra error on FindByPhone",
			inputID:    userID,
			inputPhone: otherPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByPhone", mock.Anything, otherPhone).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUserPhone,
		},
		{
			testName:   "fail - repo update error",
			inputID:    userID,
			inputPhone: otherPhone,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("FindByPhone", mock.Anything, otherPhone).Return(nil, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(dbErr)
			},
			expectedError: domain.ErrUpdatingUserPhone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			useCase := application.NewUpdateUserPhoneUseCase(mockUserRepository)
			err := useCase.Execute(context.Background(), tc.inputID, tc.inputPhone)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if tc.expectedError == domain.ErrUpdatingUserPhone && errors.Unwrap(err) != nil {
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
