package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/posts/domain"
	"jdgonzalez907/saas-api/internal/posts/infrastructure/controllers"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestCreatePostController_Handle(t *testing.T) {
	validBody := controllers.CreatePostRequest{
		ContentInformationDTO: domain.ContentInformationDTO{
			Title: "Post Title",
			Content: []domain.BlockDTO{
				{
					Type: "title",
					Text: "Title Block",
				},
			},
		},
		Status: "draft",
	}

	testCases := []struct {
		testName       string
		authUserID     any
		requestBody    any
		setupMock      func(m *mockApp.MockCreatePostUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:   "success - post created",
			authUserID: int64(1),
			requestBody: validBody,
			setupMock: func(m *mockApp.MockCreatePostUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.MatchedBy(func(p *domain.Post) bool {
					return p.ContentInformation().Title() == "Post Title" &&
						p.Status() == domain.StatusDraft &&
						p.AuthorID() == int64(1)
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"title":"Post Title"`,
		},
		{
			testName:       "fail - unauthenticated",
			authUserID:     nil,
			requestBody:    validBody,
			setupMock:      func(_ *mockApp.MockCreatePostUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   sharedHttp.ErrUnauthenticated.Error(),
		},
		{
			testName:       "fail - invalid request body",
			authUserID:     int64(1),
			requestBody:    "{invalid json}",
			setupMock:      func(_ *mockApp.MockCreatePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   sharedHttp.ErrInvalidRequestBody.Error(),
		},
		{
			testName:   "fail - empty post title",
			authUserID: int64(1),
			requestBody: controllers.CreatePostRequest{
				ContentInformationDTO: domain.ContentInformationDTO{
					Title: "",
				},
				Status: "draft",
			},
			setupMock:      func(_ *mockApp.MockCreatePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrEmptyPostTitle.Error(),
		},
		{
			testName:   "fail - invalid status",
			authUserID: int64(1),
			requestBody: controllers.CreatePostRequest{
				ContentInformationDTO: validBody.ContentInformationDTO,
				Status:                "invalid-status",
			},
			setupMock:      func(_ *mockApp.MockCreatePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPostStatus.Error(),
		},
		{
			testName:    "fail - invalid author ID",
			authUserID:  int64(0),
			requestBody: validBody,
			setupMock:   func(_ *mockApp.MockCreatePostUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidAuthorID.Error(),
		},
		{
			testName:   "fail - usecase execution error",
			authUserID: int64(1),
			requestBody: validBody,
			setupMock: func(m *mockApp.MockCreatePostUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything).Return(errors.New("db insert failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockCreatePostUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewCreatePostController(mockUseCase)

			var buf bytes.Buffer
			if tc.requestBody != nil {
				if s, ok := tc.requestBody.(string); ok {
					buf.WriteString(s)
				} else {
					err := json.NewEncoder(&buf).Encode(tc.requestBody)
					assert.NoError(t, err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/posts", &buf)
			if tc.authUserID != nil {
				req = req.WithContext(sharedHttp.WithUserID(req.Context(), tc.authUserID.(int64)))
			}

			rec := httptest.NewRecorder()
			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}
