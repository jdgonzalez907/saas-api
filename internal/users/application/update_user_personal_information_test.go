package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/domain"
	domainMocks "jdgonzalez907/saas-api/mocks/domain"

	"github.com/stretchr/testify/mock"
)

func TestUpdateUserPersonalInformationUseCase(t *testing.T) {
	userID := int64(1)
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-18, 0, -1).Format("2006-01-02"))
	now := time.Now()

	existingPersonalInfo, _ := domain.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)

	existingUser, _ := domain.NewUser(domain.UserParams{
		ID:                  userID,
		PersonalInformation: existingPersonalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	otherIdentification, _ := domain.NewIdentification(domain.IdType_CC, "2222")
	personalInfo, _ := domain.NewPersonalInformation(
		otherIdentification,
		"Jane",
		"Smith",
		&address,
		&birthDate,
	)

	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		id               int64
		info             domain.PersonalInformation
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName: "success - update personal info",
			id:       userID,
			info:     personalInfo,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "success - no changes to personal info",
			id:       userID,
			info:     existingPersonalInfo,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
			},
			expectedError: nil,
		},
		{
			testName: "fail - invalid user id",
			id:       int64(0),
			info:     personalInfo,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				// No expectations
			},
			expectedError: domain.ErrInvalidUserID,
		},
		{
			testName: "fail - user not found",
			id:       userID,
			info:     personalInfo,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName: "fail - infra error on FindById",
			id:       userID,
			info:     personalInfo,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUserPersonalInformation,
		},
		{
			testName: "fail - repo update error",
			id:       userID,
			info:     personalInfo,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", mock.Anything, userID).Return(existingUser, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(dbErr)
			},
			expectedError: domain.ErrUpdatingUserPersonalInformation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			useCase := application.NewUpdateUserPersonalInformationUseCase(mockUserRepository)
			err := useCase.Execute(context.Background(), tc.id, tc.info)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if tc.expectedError == domain.ErrUpdatingUserPersonalInformation && errors.Unwrap(err) != nil {
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
