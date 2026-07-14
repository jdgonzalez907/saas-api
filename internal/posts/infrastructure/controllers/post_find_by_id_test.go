package controllers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/posts/domain"
	"jdgonzalez907/saas-api/internal/posts/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestFindPostByIDController_Handle(t *testing.T) {
	titleBlock, _ := domain.NewTitleBlock("Title")
	contentInfo, _ := domain.NewContentInformation("Post Title", []domain.Block{titleBlock})
	now := time.Now().UTC()
	validPost, err := domain.NewPost(domain.PostParams{
		ID:                 1,
		ContentInformation: contentInfo,
		Status:             domain.StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		AuthorID:           2,
		LastEditorID:       2,
	})
	assert.NoError(t, err)

	testCases := []struct {
		testName       string
		routeParamID   string
		setupMock      func(m *mockApp.MockFindPostByIDUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - post found",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindPostByIDUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1)).Return(validPost, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Post Title"`,
		},
		{
			testName:       "fail - route parameter is not an integer",
			routeParamID:   "abc",
			setupMock:      func(_ *mockApp.MockFindPostByIDUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:     "fail - post not found",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindPostByIDUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1)).Return(nil, domain.ErrPostNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrPostNotFound.Error(),
		},
		{
			testName:     "fail - usecase execution error",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindPostByIDUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1)).Return(nil, errors.New("db find failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockFindPostByIDUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewFindPostByIDController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/posts", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.routeParamID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}
