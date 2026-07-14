package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestUpdateUserPersonalInformationController_Handle(t *testing.T) {
	validBirthDate := time.Now().AddDate(-25, 0, 0)
	validBirthDateDTO := domain.BirthDateDTO(validBirthDate.Format("2006-01-02"))

	validBody := domain.PersonalInformationDTO{
		FirstName: "Jane",
		LastName:  "Smith",
		Identification: domain.IdentificationDTO{
			Type:   domain.IDTypeCC,
			Number: "987654321",
		},
		Address: &domain.AddressDTO{
			Street:  "Street 456",
			City:    "Boston",
			State:   "MA",
			Country: "USA",
		},
		BirthDate: &validBirthDateDTO,
	}

	testCases := []struct {
		testName       string
		routeParamID   string
		requestBody    any
		setupMock      func(m *mockApp.MockUpdateUserPersonalInformationUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - personal information updated",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPersonalInformationUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.MatchedBy(func(u domain.PersonalInformation) bool {
					dto := u.ToDTO()
					return dto.FirstName == "Jane" &&
						dto.LastName == "Smith" &&
						dto.Identification.Number == "987654321"
				})).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName:       "fail - route parameter is not an integer",
			routeParamID:   "abc",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - route parameter is empty",
			routeParamID:   "",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id is missing",
		},
		{
			testName:       "fail - route parameter is negative",
			routeParamID:   "-1",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - invalid json body",
			routeParamID:   "1",
			requestBody:    "{invalid json}",
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   sharedHttp.ErrInvalidRequestBody.Error(),
		},
		{
			testName:       "fail - nil request body",
			routeParamID:   "1",
			requestBody:    nil,
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   sharedHttp.ErrInvalidRequestBody.Error(),
		},
		{
			testName:     "fail - invalid identification type",
			routeParamID: "1",
			requestBody: func() domain.PersonalInformationDTO {
				b := validBody
				b.Identification.Type = "INVALID"
				return b
			}(),
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidIdentificationType.Error(),
		},
		{
			testName:     "fail - invalid address street",
			routeParamID: "1",
			requestBody: func() domain.PersonalInformationDTO {
				b := validBody
				b.Address = &domain.AddressDTO{
					Street:  "",
					City:    "Boston",
					State:   "MA",
					Country: "USA",
				}
				return b
			}(),
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidStreet.Error(),
		},
		{
			testName:     "fail - invalid birth date (too young)",
			routeParamID: "1",
			requestBody: func() domain.PersonalInformationDTO {
				b := validBody
				bd := domain.BirthDateDTO(time.Now().AddDate(-5, 0, 0).Format("2006-01-02"))
				b.BirthDate = &bd
				return b
			}(),
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrUserUnderage.Error(),
		},
		{
			testName:     "fail - invalid birth date format",
			routeParamID: "1",
			requestBody: func() domain.PersonalInformationDTO {
				b := validBody
				bd := domain.BirthDateDTO("invalid-format")
				b.BirthDate = &bd
				return b
			}(),
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidBirthDateFormat.Error(),
		},
		{
			testName:     "fail - invalid first name",
			routeParamID: "1",
			requestBody: func() domain.PersonalInformationDTO {
				b := validBody
				b.FirstName = ""
				return b
			}(),
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidFirstName.Error(),
		},
		{
			testName:     "fail - invalid last name",
			routeParamID: "1",
			requestBody: func() domain.PersonalInformationDTO {
				b := validBody
				b.LastName = ""
				return b
			}(),
			setupMock:      func(_ *mockApp.MockUpdateUserPersonalInformationUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidLastName.Error(),
		},
		{
			testName:     "fail - user not found",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPersonalInformationUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrUserNotFound.Error(),
		},
		{
			testName:     "fail - internal server error from usecase",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPersonalInformationUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(errors.New("db update failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockUpdateUserPersonalInformationUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewUpdateUserPersonalInformationController(mockUseCase)

			var req *http.Request
			if tc.requestBody == nil {
				req = httptest.NewRequest(http.MethodPut, "/users", nil)
				req.Body = nil
			} else {
				var buf bytes.Buffer
				if s, ok := tc.requestBody.(string); ok {
					buf.WriteString(s)
				} else {
					err := json.NewEncoder(&buf).Encode(tc.requestBody)
					assert.NoError(t, err)
				}
				req = httptest.NewRequest(http.MethodPut, "/users", &buf)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.routeParamID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				var jsonResponse map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &jsonResponse)
				assert.NoError(t, err)
				assert.Contains(t, jsonResponse["message"], tc.expectedBody)
			}
		})
	}
}
