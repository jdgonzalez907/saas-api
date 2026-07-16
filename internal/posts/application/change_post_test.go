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

func TestChangePostUseCase(t *testing.T) {
	titleBlock, _ := domain.NewTitleBlock("Title")
	contentInfo, _ := domain.NewContentInformation("Post Title", []domain.Block{titleBlock})
	contentInfoUpdated, _ := domain.NewContentInformation("Updated Title", []domain.Block{titleBlock})
	now := time.Now().UTC()

	post, _ := domain.NewPost(1, contentInfo, domain.StatusDraft, now, now, 10, 10, nil)

	dbErr := errors.New("database connection error")

	testCases := []struct {
		name             string
		postID           int64
		contentInfo      domain.ContentInformation
		status           domain.PostStatus
		lastEditorID     int64
		mockExpectations func(*domainMocks.MockPostRepository)
		expectedError    error
	}{
		{
			name:         "success - change post",
			postID:       1,
			contentInfo:  contentInfoUpdated,
			status:       domain.StatusPublished,
			lastEditorID: 11,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(p *domain.Post) bool {
					return p.ID() == 1 &&
						p.ContentInformation().Equals(contentInfoUpdated) &&
						p.Status() == domain.StatusPublished &&
						p.LastEditorID() == 11
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:         "fail - post not found",
			postID:       1,
			contentInfo:  contentInfoUpdated,
			status:       domain.StatusPublished,
			lastEditorID: 11,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			expectedError: domain.ErrPostNotFound,
		},
		{
			name:         "fail - repository find error",
			postID:       1,
			contentInfo:  contentInfoUpdated,
			status:       domain.StatusPublished,
			lastEditorID: 11,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(nil, dbErr)
			},
			expectedError: domain.ErrChangingPost,
		},
		{
			name:         "fail - domain validation error (invalid editor ID)",
			postID:       1,
			contentInfo:  contentInfoUpdated,
			status:       domain.StatusPublished,
			lastEditorID: 0,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
			},
			expectedError: domain.ErrInvalidLastEditorID,
		},
		{
			name:         "fail - repository update error",
			postID:       1,
			contentInfo:  contentInfoUpdated,
			status:       domain.StatusPublished,
			lastEditorID: 11,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(dbErr)
			},
			expectedError: domain.ErrChangingPost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostRepository := new(domainMocks.MockPostRepository)
			tc.mockExpectations(mockPostRepository)

			useCase := application.NewChangePostUseCase(mockPostRepository)
			_, err := useCase.Execute(context.Background(), tc.postID, tc.contentInfo, tc.status, tc.lastEditorID)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrChangingPost && errors.Unwrap(err) != nil) {
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
