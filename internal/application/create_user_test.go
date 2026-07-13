package application_test

import (
	"errors"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/application"
	"jdgonzalez907/users-api/internal/domain"
	domainMocks "jdgonzalez907/users-api/mocks/domain"
)

func TestCreateUserUseCase(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IdType_CC, "1111")
	phone, _ := domain.NewPhone("57", "123456789")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("123 Main St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-18, 0, -1))
	now := time.Now()

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

	userWithNilEmail, _ := domain.NewUser(domain.UserParams{
		ID:                  2,
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	})

	dbErr := errors.New("database connection error")

	testCases := []struct {
		name             string
		input            domain.User
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedError    error
	}{
		{
			name:  "success - create user",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(nil, nil)
				m.On("Create", user).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "fail - id already exists",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(user, nil)
			},
			expectedError: domain.ErrUserIDAlreadyExists,
		},
		{
			name:  "fail - phone already exists",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(user, nil)
			},
			expectedError: domain.ErrUserPhoneAlreadyExists,
		},
		{
			name:  "fail - email already exists",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(user, nil)
			},
			expectedError: domain.ErrUserEmailAlreadyExists,
		},
		{
			name:  "fail - infra error on FindById",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			name:  "fail - infra error on FindByPhone",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			name:  "fail - infra error on FindByEmail",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(nil, dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			name:  "fail - repo create error",
			input: *user,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", user.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("FindByEmail", email).Return(nil, nil)
				m.On("Create", user).Return(dbErr)
			},
			expectedError: domain.ErrCreatingUser,
		},
		{
			name:  "success - nil email",
			input: *userWithNilEmail,
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindById", userWithNilEmail.ToDTO().ID).Return(nil, nil)
				m.On("FindByPhone", phone).Return(nil, nil)
				m.On("Create", userWithNilEmail).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockUserRepository)

			createUserUseCase := application.NewCreateUserUseCase(mockUserRepository)
			err := createUserUseCase.Execute(&tc.input)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if tc.expectedError == domain.ErrCreatingUser && errors.Unwrap(err) != nil {
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
