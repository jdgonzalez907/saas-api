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

func TestCreateUserUseCase(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-18, 0, -1))
	user, _ := domain.NewUser(
		uuid.NewString(),
		identification,
		"John",
		"Doe",
		phone,
		&email,
		&address,
		&birthDate,
	)

	userWithNilEmail, _ := domain.NewUser(
		uuid.NewString(),
		identification,
		"John",
		"Doe",
		phone,
		nil,
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
			testName: "create user with success",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(nil, nil)
				m.On("Create", user).Return(nil)
			},
			expectedError: nil,
		},
		{
			testName: "create user fails - id already exists",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(user, nil)
			},
			expectedError: domain.ErrUserIDAlreadyExists,
		},
		{
			testName: "create user fails - phone already exists",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(user, nil)
			},
			expectedError: domain.ErrUserPhoneAlreadyExists,
		},
		{
			testName: "create user fails - email already exists",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(user, nil)
			},
			expectedError: domain.ErrUserEmailAlreadyExists,
		},
		{
			testName: "create user fails - infra error on FindById",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			testName: "create user fails - infra error on FindByPhone",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			testName: "create user fails - infra error on FindByEmail",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(nil, dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			testName: "create user fails - repo create error",
			input:    *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(nil, nil)
				m.On("Create", user).Return(dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			testName: "create user with nil email success",
			input:    *userWithNilEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userWithNilEmail.ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("Create", userWithNilEmail).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			testCase.mockExpectations(mockUserRepository)

			createUserUseCase := application.NewCreateUserUseCase(mockUserRepository)
			err := createUserUseCase.Execute(&testCase.input)

			if testCase.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", testCase.expectedError)
				}
				if !errors.Is(err, testCase.expectedError) {
					if testCase.expectedError == domain.ErrCreatingUser && errors.Unwrap(err) != nil {
						// Success: it is a wrapped error containing the underlying cause
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
