package controllers_test

import (
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

func TestDeleteUserController_Handle(t *testing.T) {
	testCases := []struct {
		testName       string
		routeParamID   string
		setupMock      func(m *mockApp.MockDeleteUserUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - user deleted",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockDeleteUserUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName:       "fail - route parameter is not an integer",
			routeParamID:   "abc",
			setupMock:      func(_ *mockApp.MockDeleteUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - route parameter is empty",
			routeParamID:   "",
			setupMock:      func(_ *mockApp.MockDeleteUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id is missing",
		},
		{
			testName:       "fail - route parameter is negative",
			routeParamID:   "-3",
			setupMock:      func(_ *mockApp.MockDeleteUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:     "fail - user not found",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockDeleteUserUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1)).Return(domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrUserNotFound.Error(),
		},
		{
			testName:     "fail - internal server error from usecase",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockDeleteUserUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1)).Return(errors.New("db delete failure"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockDeleteUserUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewDeleteUserController(mockUseCase)

			req := httptest.NewRequest(http.MethodDelete, "/users", nil)
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
