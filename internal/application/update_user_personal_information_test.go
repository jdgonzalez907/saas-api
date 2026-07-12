package application_test

import (
	"errors"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
	domainMocks "jdgonzalez907/users-api/mocks/domain"

	"github.com/stretchr/testify/mock"
)

func TestUpdateUserPersonalInformationUseCase(t *testing.T) {
	userID := 1
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("123456789")
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

	otherIdentification, _ := domain.NewIdentification(domain.IdType_CC, "2222")
	userWithNewPersonalInformation, _ := domain.NewUser(domain.UserParams{
		ID:             userID,
		Identification: otherIdentification,
		FirstName:      "Jane",
		LastName:       "Smith",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	// User that will fail validation in WithPersonalInformation method (e.g. empty firstName)
	invalidUser, _ := domain.NewUser(domain.UserParams{
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
	rf := reflect.ValueOf(invalidUser).Elem().FieldByName("firstName")
	reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().SetString("")

	dbErr := errors.New("database connection error")

	testCases := []struct {
		testName         string
		input            domain.User
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			testName: "success - update personal info",
			input:    *userWithNewPersonalInformation,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("Update", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "fail - user not found",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(nil, nil)
			},
			expectedError: domain.ErrUserNotFound,
		},
		{
			testName: "fail - infra error on FindById",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(nil, dbErr)
			},
			expectedError: domain.ErrUpdatingUserPersonalInformation,
		},
		{
			testName: "fail - repo update error",
			input:    *existingUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
				m.On("Update", mock.Anything).Return(dbErr)
			},
			expectedError: domain.ErrUpdatingUserPersonalInformation,
		},
		{
			testName: "fail - invalid firstName",
			input:    *invalidUser,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userID).Return(existingUser, nil)
			},
			expectedError: domain.ErrUpdatingUserPersonalInformation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			testCase.mockExpectations(mockUserRepository)

			useCase := application.NewUpdateUserPersonalInformationUseCase(mockUserRepository)
			err := useCase.Execute(&testCase.input)

			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", testCase.expectedError)
				}
				if !errors.Is(err, testCase.expectedError) {
					if testCase.expectedError == domain.ErrUpdatingUserPersonalInformation && errors.Unwrap(err) != nil {
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
