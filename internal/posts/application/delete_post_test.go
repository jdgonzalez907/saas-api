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

func TestDeletePostUseCase(t *testing.T) {
	titleBlock, _ := domain.NewTitleBlock("Title")
	contentInfo, _ := domain.NewContentInformation("Post Title", []domain.Block{titleBlock})
	now := time.Now().UTC()

	post, _ := domain.NewPost(domain.PostParams{
		ID:                 1,
		ContentInformation: contentInfo,
		Status:             domain.StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		AuthorID:           10,
		LastEditorID:       10,
	})

	dbErr := errors.New("database connection error")

	testCases := []struct {
		name             string
		postID           int64
		deletedByID      int64
		mockExpectations func(*domainMocks.MockPostRepository)
		expectedError    error
	}{
		{
			name:        "success - delete post",
			postID:      1,
			deletedByID: 10,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				m.On("Delete", mock.Anything, int64(1), int64(10)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "fail - post not found",
			postID:      1,
			deletedByID: 10,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			expectedError: domain.ErrPostNotFound,
		},
		{
			name:        "fail - repository find error",
			postID:      1,
			deletedByID: 10,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(nil, dbErr)
			},
			expectedError: domain.ErrDeletingPost,
		},
		{
			name:        "fail - repository delete error",
			postID:      1,
			deletedByID: 10,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				m.On("Delete", mock.Anything, int64(1), int64(10)).Return(dbErr)
			},
			expectedError: domain.ErrDeletingPost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostRepository := new(domainMocks.MockPostRepository)
			tc.mockExpectations(mockPostRepository)

			useCase := application.NewDeletePostUseCase(mockPostRepository)
			err := useCase.Execute(context.Background(), tc.postID, tc.deletedByID)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrDeletingPost && errors.Unwrap(err) != nil) {
						t.Errorf("expected error to wrap or be %v, got %v", tc.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}
