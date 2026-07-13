package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
	"jdgonzalez907/users-api/internal/infrastructure/controllers"
	mockApp "jdgonzalez907/users-api/mocks/application"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserController_Handle(t *testing.T) {
	validBirthDate := time.Now().AddDate(-25, 0, 0)
	validBirthDateDTO := domain.BirthDateDTO(validBirthDate.Format("2006-01-02"))

	validEmail := domain.EmailDTO("john.doe@example.com")

	validBody := domain.UserDTO{
		PersonalInformationDTO: domain.PersonalInformationDTO{
			FirstName: "John",
			LastName:  "Doe",
			Identification: domain.IdentificationDTO{
				Type:   domain.IdType_CC,
				Number: "123456789",
			},
			Address: &domain.AddressDTO{
				Street:  "Street 123",
				City:    "New York",
				State:   "NY",
				Country: "USA",
			},
			BirthDate: &validBirthDateDTO,
		},
		Phone: domain.PhoneDTO{
			CountryCode: "57",
			Number:      "3112223344",
		},
		Email: &validEmail,
	}

	testCases := []struct {
		testName       string
		requestBody    any
		setupMock      func(m *mockApp.MockCreateUserUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:    "success - user created",
			requestBody: validBody,
			setupMock: func(m *mockApp.MockCreateUserUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.FirstName() == "John" &&
						u.LastName() == "Doe" &&
						u.Phone().Number() == "3112223344" &&
						u.Phone().CountryCode() == "57"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			testName:       "fail - invalid json body",
			requestBody:    "{invalid json}",
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   controllers.ErrInvalidRequestBody.Error(),
		},
		{
			testName:       "fail - nil/empty body",
			requestBody:    nil,
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   controllers.ErrInvalidRequestBody.Error(),
		},
		{
			testName: "fail - invalid identification type",
			requestBody: func() domain.UserDTO {
				b := validBody
				b.Identification.Type = "INVALID"
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidIdentificationType.Error(),
		},
		{
			testName: "fail - invalid phone number",
			requestBody: func() domain.UserDTO {
				b := validBody
				b.Phone.Number = ""
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPhone.Error(),
		},
		{
			testName: "fail - invalid email",
			requestBody: func() domain.UserDTO {
				b := validBody
				badEmail := domain.EmailDTO("invalid-email")
				b.Email = &badEmail
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidEmail.Error(),
		},
		{
			testName: "fail - invalid address street",
			requestBody: func() domain.UserDTO {
				b := validBody
				b.Address = &domain.AddressDTO{
					Street:  "",
					City:    "New York",
					State:   "NY",
					Country: "USA",
				}
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidStreet.Error(),
		},
		{
			testName: "fail - invalid birth date (too young)",
			requestBody: func() domain.UserDTO {
				b := validBody
				bd := domain.BirthDateDTO(time.Now().AddDate(-5, 0, 0).Format("2006-01-02"))
				b.BirthDate = &bd
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrUserUnderage.Error(),
		},
		{
			testName: "fail - invalid birth date format",
			requestBody: func() domain.UserDTO {
				b := validBody
				bd := domain.BirthDateDTO("invalid-format")
				b.BirthDate = &bd
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidBirthDateFormat.Error(),
		},
		{
			testName: "fail - invalid user first name",
			requestBody: func() domain.UserDTO {
				b := validBody
				b.FirstName = ""
				return b
			}(),
			setupMock:      func(m *mockApp.MockCreateUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidFirstName.Error(),
		},
		{
			testName:    "fail - user with phone already exists (conflict)",
			requestBody: validBody,
			setupMock: func(m *mockApp.MockCreateUserUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything).Return(domain.ErrUserPhoneAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   domain.ErrUserPhoneAlreadyExists.Error(),
		},
		{
			testName:    "fail - internal server error during creation",
			requestBody: validBody,
			setupMock: func(m *mockApp.MockCreateUserUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything).Return(errors.New("db insert failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockCreateUserUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewCreateUserController(mockUseCase)

			var req *http.Request
			if tc.requestBody == nil {
				req = httptest.NewRequest(http.MethodPost, "/users", nil)
				req.Body = nil
			} else {
				var buf bytes.Buffer
				if s, ok := tc.requestBody.(string); ok {
					buf.WriteString(s)
				} else {
					err := json.NewEncoder(&buf).Encode(tc.requestBody)
					assert.NoError(t, err)
				}
				req = httptest.NewRequest(http.MethodPost, "/users", &buf)
			}

			rec := httptest.NewRecorder()
			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedStatus == http.StatusCreated {
				var res domain.UserDTO
				err := json.Unmarshal(rec.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, "John", res.FirstName)
				assert.Equal(t, "Doe", res.LastName)
			} else if tc.expectedBody != "" {
				var jsonResponse map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &jsonResponse)
				assert.NoError(t, err)
				assert.Contains(t, jsonResponse["message"], tc.expectedBody)
			}
		})
	}
}
