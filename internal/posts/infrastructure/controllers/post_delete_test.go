package controllers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/posts/domain"
	"jdgonzalez907/saas-api/internal/posts/infrastructure/controllers"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestDeletePostController_Handle(t *testing.T) {
	testCases := []struct {
		testName       string
		authUserID     any
		routeParamID   string
		setupMock      func(m *mockApp.MockDeletePostUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - post deleted",
			authUserID:   int64(3),
			routeParamID: "1",
			setupMock: func(m *mockApp.MockDeletePostUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), int64(3)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName:       "fail - unauthenticated",
			authUserID:     nil,
			routeParamID:   "1",
			setupMock:      func(_ *mockApp.MockDeletePostUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   sharedHttp.ErrUnauthenticated.Error(),
		},
		{
			testName:       "fail - route parameter is not an integer",
			authUserID:     int64(3),
			routeParamID:   "abc",
			setupMock:      func(_ *mockApp.MockDeletePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:     "fail - post not found",
			authUserID:   int64(3),
			routeParamID: "1",
			setupMock: func(m *mockApp.MockDeletePostUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), int64(3)).Return(domain.ErrPostNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrPostNotFound.Error(),
		},
		{
			testName:     "fail - usecase execution error",
			authUserID:   int64(3),
			routeParamID: "1",
			setupMock: func(m *mockApp.MockDeletePostUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), int64(3)).Return(errors.New("db delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockDeletePostUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewDeletePostController(mockUseCase)

			req := httptest.NewRequest(http.MethodDelete, "/posts", nil)
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
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}
