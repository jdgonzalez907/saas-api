package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestUpdateUserPhoneController_Handle(t *testing.T) {
	validBody := domain.PhoneDTO{CountryCode: "57", Number: "987654321"}

	testCases := []struct {
		testName       string
		authUserID     any
		routeParamID   string
		requestBody    any
		setupMock      func(m *mockApp.MockUpdateUserPhoneUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - phone updated",
			authUserID:   int64(1),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPhoneUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.MatchedBy(func(p domain.Phone) bool {
					return p.CountryCode() == "57" && p.Number() == "987654321"
				})).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName:       "fail - unauthenticated",
			authUserID:     nil,
			routeParamID:   "1",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockUpdateUserPhoneUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   sharedHttp.ErrUnauthenticated.Error(),
		},
		{
			testName:       "fail - route parameter is not an integer",
			authUserID:     int64(1),
			routeParamID:   "abc",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockUpdateUserPhoneUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - invalid json body",
			authUserID:     int64(1),
			routeParamID:   "1",
			requestBody:    "{invalid json}",
			setupMock:      func(_ *mockApp.MockUpdateUserPhoneUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   sharedHttp.ErrInvalidRequestBody.Error(),
		},
		{
			testName:       "fail - nil request body",
			authUserID:     int64(1),
			routeParamID:   "1",
			requestBody:    nil,
			setupMock:      func(_ *mockApp.MockUpdateUserPhoneUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   sharedHttp.ErrInvalidRequestBody.Error(),
		},
		{
			testName:     "fail - invalid phone country code",
			authUserID:   int64(1),
			routeParamID: "1",
			requestBody: domain.PhoneDTO{
				CountryCode: "", Number: "987654321",
			},
			setupMock:      func(_ *mockApp.MockUpdateUserPhoneUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPhone.Error(),
		},
		{
			testName:     "fail - invalid phone number",
			authUserID:   int64(1),
			routeParamID: "1",
			requestBody: domain.PhoneDTO{
				CountryCode: "57", Number: "",
			},
			setupMock:      func(_ *mockApp.MockUpdateUserPhoneUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPhone.Error(),
		},
		{
			testName:     "fail - user not found",
			authUserID:   int64(1),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPhoneUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrUserNotFound.Error(),
		},
		{
			testName:     "fail - phone already exists",
			authUserID:   int64(1),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPhoneUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(domain.ErrUserPhoneAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   domain.ErrUserPhoneAlreadyExists.Error(),
		},
		{
			testName:     "fail - internal server error from usecase",
			authUserID:   int64(1),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserPhoneUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(errors.New("db failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockUpdateUserPhoneUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewUpdateUserPhoneController(mockUseCase)

			var req *http.Request
			if tc.requestBody == nil {
				req = httptest.NewRequest(http.MethodPut, "/users/phone", nil)
				req.Body = nil
			} else {
				var buf bytes.Buffer
				if s, ok := tc.requestBody.(string); ok {
					buf.WriteString(s)
				} else {
					err := json.NewEncoder(&buf).Encode(tc.requestBody)
					assert.NoError(t, err)
				}
				req = httptest.NewRequest(http.MethodPut, "/users/phone", &buf)
			}

			if tc.authUserID != nil {
				req.Header.Set("Authorization", strconv.FormatInt(tc.authUserID.(int64), 10))
			}
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.routeParamID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			handler := sharedHttp.Protected(controller.Handle)
			handler.ServeHTTP(rec, req)

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
