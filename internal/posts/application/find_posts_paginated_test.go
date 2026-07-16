package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/posts/application"
	"jdgonzalez907/saas-api/internal/posts/domain"
	domainMocks "jdgonzalez907/saas-api/mocks/domain"
)

func TestFindPostsPaginatedUseCase(t *testing.T) {
	titleBlock, _ := domain.NewTitleBlock("Title")
	contentInfo, _ := domain.NewContentInformation("Post Title", []domain.Block{titleBlock})
	now := time.Now().UTC()

	post, _ := domain.NewPost(1, contentInfo, domain.StatusPublished, now, now, 10, 10, &now)

	pagination, _ := domain.NewPagination(nil, nil, nil)
	dbErr := errors.New("database connection error")

	testCases := []struct {
		name             string
		status           domain.PostStatus
		pagination       domain.Pagination
		mockExpectations func(*domainMocks.MockPostRepository)
		expectedResult   domain.PaginatedPosts
		expectedError    error
	}{
		{
			name:       "success - paginated posts not empty",
			status:     domain.StatusPublished,
			pagination: pagination,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindAll", mock.Anything, domain.StatusPublished, pagination).Return([]*domain.Post{post}, nil)
			},
			expectedResult: domain.NewPaginatedPosts([]*domain.Post{post}, &now, func() *int64 { i := int64(1); return &i }()),
			expectedError:  nil,
		},
		{
			name:       "success - paginated posts empty",
			status:     domain.StatusPublished,
			pagination: pagination,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindAll", mock.Anything, domain.StatusPublished, pagination).Return([]*domain.Post{}, nil)
			},
			expectedResult: domain.NewPaginatedPosts([]*domain.Post{}, nil, nil),
			expectedError:  nil,
		},
		{
			name:       "fail - repository error",
			status:     domain.StatusPublished,
			pagination: pagination,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindAll", mock.Anything, domain.StatusPublished, pagination).Return(nil, dbErr)
			},
			expectedResult: domain.PaginatedPosts{},
			expectedError:  domain.ErrFindingPosts,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostRepository := new(domainMocks.MockPostRepository)
			tc.mockExpectations(mockPostRepository)

			useCase := application.NewFindPostsPaginatedUseCase(mockPostRepository)
			res, err := useCase.Execute(context.Background(), tc.status, tc.pagination)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrFindingPosts && errors.Unwrap(err) != nil) {
						t.Errorf("expected error to wrap or be %v, got %v", tc.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if len(res.Posts()) != len(tc.expectedResult.Posts()) {
					t.Errorf("expected %d posts, got %d", len(tc.expectedResult.Posts()), len(res.Posts()))
				}
				if (res.NextPublishedAt() == nil) != (tc.expectedResult.NextPublishedAt() == nil) {
					t.Error("mismatched next published at presence")
				}
				if res.NextPublishedAt() != nil && !res.NextPublishedAt().Equal(*tc.expectedResult.NextPublishedAt()) {
					t.Errorf("expected next published at %v, got %v", tc.expectedResult.NextPublishedAt(), res.NextPublishedAt())
				}
				if (res.NextID() == nil) != (tc.expectedResult.NextID() == nil) {
					t.Error("mismatched next id presence")
				}
				if res.NextID() != nil && *res.NextID() != *tc.expectedResult.NextID() {
					t.Errorf("expected next id %d, got %d", *tc.expectedResult.NextID(), *res.NextID())
				}
			}
		})
	}
}
