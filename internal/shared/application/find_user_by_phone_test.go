package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"jdgonzalez907/saas-api/internal/shared/application"
	"jdgonzalez907/saas-api/internal/shared/domain"
	domainUser "jdgonzalez907/saas-api/internal/users/domain"
	domainMocks "jdgonzalez907/saas-api/mocks/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	identification, _ := domainUser.NewIdentification(domainUser.IDTypeCC, "123456789")
	phone, _ := domainUser.NewPhone("57", "3112223344")
	address, _ := domainUser.NewAddress("Street 1", "City", "State", "Country", nil, nil)
	email, _ := domainUser.NewEmail("test@example.com")
	birthDate, _ := domainUser.NewBirthDate(time.Now().AddDate(-20, 0, 0).Format("2006-01-02"))

	personalInfo, _ := domainUser.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)
	now := time.Now().UTC()

	userFound, err := domainUser.NewUser(domainUser.UserParams{
		ID:                  1,
		PersonalInformation: personalInfo,
		Phone:               phone,
		Email:               &email,
		CreatedAt:           now,
		UpdatedAt:           now,
	})
	assert.NoError(t, err)

	var emailVal string
	if userFound.Email() != nil {
		emailVal = userFound.Email().Value()
	}

	validUser := domain.NewUser(userFound.ID(),
		userFound.FullName(),
		userFound.Phone().CountryCode(),
		userFound.Phone().Number(),
		&emailVal,
	)

	testCases := []struct {
		testName         string
		inputCountryCode string
		inputNumber      string
		mockExpectations func(*domainMocks.MockUserRepository)
		expectedResult   *domain.User
		expectedError    error
	}{
		{
			testName:         "success - find user by phone",
			inputCountryCode: "57",
			inputNumber:      "3112223344",
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByPhone", mock.Anything, phone).Return(userFound, nil)
			},
			expectedResult: validUser,
			expectedError:  nil,
		},
		{
			testName:         "fail - user not found",
			inputCountryCode: "57",
			inputNumber:      "3112223344",
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByPhone", mock.Anything, phone).Return(nil, nil)
			},
			expectedResult: nil,
			expectedError:  domainUser.ErrUserNotFound,
		},
		{
			testName:         "fail - internal server error",
			inputCountryCode: "57",
			inputNumber:      "3112223344",
			mockExpectations: func(m *domainMocks.MockUserRepository) {
				m.On("FindByPhone", mock.Anything, phone).Return(nil, errors.New("internal server error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("internal server error"),
		},
		{
			testName:         "fail - invalid phone number",
			inputCountryCode: "",
			inputNumber:      "1234567890",
			mockExpectations: func(m *domainMocks.MockUserRepository) {},
			expectedResult:   nil,
			expectedError:    domainUser.ErrInvalidPhone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepo := new(domainMocks.MockUserRepository)
			tc.mockExpectations(mockRepo)

			app := application.NewFindUserByPhoneUseCase(mockRepo)
			result, err := app.Execute(context.Background(), tc.inputCountryCode, tc.inputNumber)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
