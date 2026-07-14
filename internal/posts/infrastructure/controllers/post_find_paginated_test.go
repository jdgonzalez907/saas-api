package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/posts/domain"
	"jdgonzalez907/saas-api/internal/posts/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestFindPostsPaginatedController_Handle(t *testing.T) {
	titleBlock, _ := domain.NewTitleBlock("Title")
	contentInfo, _ := domain.NewContentInformation("Post Title", []domain.Block{titleBlock})
	now := time.Now().UTC()
	post, err := domain.NewPost(domain.PostParams{
		ID:                 1,
		ContentInformation: contentInfo,
		Status:             domain.StatusPublished,
		CreatedAt:          now,
		UpdatedAt:          now,
		AuthorID:           2,
		LastEditorID:       2,
		PublishedAt:        &now,
	})
	assert.NoError(t, err)

	nextID := int64(1)
	paginatedPosts := domain.NewPaginatedPosts([]*domain.Post{post}, &now, &nextID)

	testCases := []struct {
		testName       string
		urlQuery       string
		setupMock      func(m *mockApp.MockFindPostsPaginatedUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "success - default pagination",
			urlQuery: "?status=published",
			setupMock: func(m *mockApp.MockFindPostsPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything, domain.StatusPublished, mock.MatchedBy(func(p domain.Pagination) bool {
					return p.LastID() == nil && p.LastPublishedAt() == nil && p.Limit() == 10
				})).Return(paginatedPosts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"posts"`,
		},
		{
			testName: "success - with cursor and limit",
			urlQuery: "?status=published&limit=25&lastID=5&lastPublishedAt=2026-07-14T19:00:00Z",
			setupMock: func(m *mockApp.MockFindPostsPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything, domain.StatusPublished, mock.MatchedBy(func(p domain.Pagination) bool {
					return p.LastID() != nil && *p.LastID() == 5 &&
						p.LastPublishedAt() != nil && p.Limit() == 25
				})).Return(paginatedPosts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"posts"`,
		},
		{
			testName:       "fail - invalid status",
			urlQuery:       "?status=invalid-status",
			setupMock:      func(_ *mockApp.MockFindPostsPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPostStatus.Error(),
		},
		{
			testName:       "fail - invalid limit",
			urlQuery:       "?status=published&limit=abc",
			setupMock:      func(_ *mockApp.MockFindPostsPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter limit must be a positive integer",
		},
		{
			testName:       "fail - invalid lastID",
			urlQuery:       "?status=published&lastID=abc",
			setupMock:      func(_ *mockApp.MockFindPostsPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter lastID must be a positive integer",
		},
		{
			testName:       "fail - invalid lastPublishedAt",
			urlQuery:       "?status=published&lastPublishedAt=abc",
			setupMock:      func(_ *mockApp.MockFindPostsPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter lastPublishedAt must be a valid RFC3339 timestamp",
		},
		{
			testName:       "fail - pagination validation error (mismatched cursor params)",
			urlQuery:       "?status=published&lastID=5",
			setupMock:      func(_ *mockApp.MockFindPostsPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPaginationCursor.Error(),
		},
		{
			testName: "fail - usecase execution error",
			urlQuery: "?status=published",
			setupMock: func(m *mockApp.MockFindPostsPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(domain.PaginatedPosts{}, errors.New("db find failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockFindPostsPaginatedUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewFindPostsPaginatedController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/posts"+tc.urlQuery, nil)
			rec := httptest.NewRecorder()

			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}
