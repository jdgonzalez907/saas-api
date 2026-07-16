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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/posts/domain"
	"jdgonzalez907/saas-api/internal/posts/infrastructure/controllers"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestChangePostController_Handle(t *testing.T) {
	titleBlock, _ := domain.NewTitleBlock("Title")
	contentInfo, _ := domain.NewContentInformation("Post Title", []domain.Block{titleBlock})
	now := time.Now().UTC()
	validPost, err := domain.NewPost(1, contentInfo, domain.StatusPublished, now, now, 2, &now)
	assert.NoError(t, err)

	validBody := controllers.ChangePostRequest{
		ContentInformationDTO: domain.ContentInformationDTO{
			Title: "Post Title",
			Content: []domain.BlockDTO{
				{
					Type: "title",
					Text: "Title Block",
				},
			},
		},
		Status: "published",
	}

	testCases := []struct {
		testName       string
		authUserID     any
		routeParamID   string
		requestBody    any
		setupMock      func(m *mockApp.MockChangePostUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - post changed",
			authUserID:   int64(3),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockChangePostUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), mock.MatchedBy(func(c domain.ContentInformation) bool {
					return c.Title() == "Post Title"
				}), domain.StatusPublished, int64(3)).Return(validPost, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Post Title"`,
		},
		{
			testName:       "fail - unauthenticated",
			authUserID:     nil,
			routeParamID:   "1",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockChangePostUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   sharedHttp.ErrUnauthenticated.Error(),
		},
		{
			testName:       "fail - route parameter is not an integer",
			authUserID:     int64(3),
			routeParamID:   "abc",
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockChangePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - invalid request body",
			authUserID:     int64(3),
			routeParamID:   "1",
			requestBody:    "{invalid json}",
			setupMock:      func(_ *mockApp.MockChangePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   sharedHttp.ErrInvalidRequestBody.Error(),
		},
		{
			testName:     "fail - empty post title",
			authUserID:   int64(3),
			routeParamID: "1",
			requestBody: controllers.ChangePostRequest{
				ContentInformationDTO: domain.ContentInformationDTO{
					Title: "",
				},
				Status: "published",
			},
			setupMock:      func(_ *mockApp.MockChangePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrEmptyPostTitle.Error(),
		},
		{
			testName:     "fail - invalid status",
			authUserID:   int64(3),
			routeParamID: "1",
			requestBody: controllers.ChangePostRequest{
				ContentInformationDTO: validBody.ContentInformationDTO,
				Status:                "invalid-status",
			},
			setupMock:      func(_ *mockApp.MockChangePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPostStatus.Error(),
		},
		{
			testName:     "fail - post not found",
			authUserID:   int64(3),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockChangePostUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, domain.ErrPostNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrPostNotFound.Error(),
		},
		{
			testName:     "fail - usecase execution error",
			authUserID:   int64(3),
			routeParamID: "1",
			requestBody:  validBody,
			setupMock: func(m *mockApp.MockChangePostUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db update failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockChangePostUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewChangePostController(mockUseCase)

			var buf bytes.Buffer
			if tc.requestBody != nil {
				if s, ok := tc.requestBody.(string); ok {
					buf.WriteString(s)
				} else {
					err := json.NewEncoder(&buf).Encode(tc.requestBody)
					assert.NoError(t, err)
				}
			}

			req := httptest.NewRequest(http.MethodPut, "/posts", &buf)
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
