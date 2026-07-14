package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestUpdateUserEmailController_Handle(t *testing.T) {
	newEmail := "new@example.com"
	emptyEmail := ""
	invalidEmail := "invalid-email"

	validBody := controllers.UpdateEmailRequest{
		Email: &newEmail,
	}

	testCases := []struct {
		testName       string
		routeParamID   string
		requestBody    any
		setupMock      func(m *mockApp.MockUpdateUserEmailUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - email updated",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserEmailUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.MatchedBy(func(e *domain.Email) bool {
					return e != nil && e.Value() == "new@example.com"
				})).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName:     "success - email set to nil",
			routeParamID: "1",
			requestBody: controllers.UpdateEmailRequest{
				Email: &emptyEmail,
			},
			setupMock: func(m *mockApp.MockUpdateUserEmailUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), (*domain.Email)(nil)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName:       "fail - route parameter is not an integer",
			routeParamID:   "abc",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockUpdateUserEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - invalid json body",
			routeParamID:   "1",
			requestBody:    "{invalid json}",
			setupMock:      func(_ *mockApp.MockUpdateUserEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   controllers.ErrInvalidRequestBody.Error(),
		},
		{
			testName:       "fail - nil request body",
			routeParamID:   "1",
			requestBody:    nil,
			setupMock:      func(_ *mockApp.MockUpdateUserEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   controllers.ErrInvalidRequestBody.Error(),
		},
		{
			testName:     "fail - invalid email format",
			routeParamID: "1",
			requestBody: controllers.UpdateEmailRequest{
				Email: &invalidEmail,
			},
			setupMock:      func(_ *mockApp.MockUpdateUserEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidEmail.Error(),
		},
		{
			testName:     "fail - user not found",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserEmailUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrUserNotFound.Error(),
		},
		{
			testName:     "fail - email already exists",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserEmailUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(domain.ErrUserEmailAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   domain.ErrUserEmailAlreadyExists.Error(),
		},
		{
			testName:     "fail - internal server error from usecase",
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockUpdateUserEmailUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.Anything).Return(errors.New("db failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockUpdateUserEmailUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewUpdateUserEmailController(mockUseCase)

			var req *http.Request
			if tc.requestBody == nil {
				req = httptest.NewRequest(http.MethodPut, "/users/email", nil)
				req.Body = nil
			} else {
				var buf bytes.Buffer
				if s, ok := tc.requestBody.(string); ok {
					buf.WriteString(s)
				} else {
					err := json.NewEncoder(&buf).Encode(tc.requestBody)
					assert.NoError(t, err)
				}
				req = httptest.NewRequest(http.MethodPut, "/users/email", &buf)
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
