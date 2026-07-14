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

func TestFindPostByIDUseCase(t *testing.T) {
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
		inputID          int64
		mockExpectations func(*domainMocks.MockPostRepository)
		expectedResult   *domain.Post
		expectedError    error
	}{
		{
			name:    "success - find post",
			inputID: 1,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
			},
			expectedResult: post,
			expectedError:  nil,
		},
		{
			name:    "fail - post not found",
			inputID: 1,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			expectedResult: nil,
			expectedError:  domain.ErrPostNotFound,
		},
		{
			name:    "fail - repository error",
			inputID: 1,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("FindByID", mock.Anything, int64(1)).Return(nil, dbErr)
			},
			expectedResult: nil,
			expectedError:  domain.ErrFindingPost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostRepository := new(domainMocks.MockPostRepository)
			tc.mockExpectations(mockPostRepository)

			useCase := application.NewFindPostByIDUseCase(mockPostRepository)
			res, err := useCase.Execute(context.Background(), tc.inputID)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrFindingPost && errors.Unwrap(err) != nil) {
						t.Errorf("expected error to wrap or be %v, got %v", tc.expectedError, err)
					}
				}
				if res != nil {
					t.Error("expected returned post to be nil on error")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if !res.Equals(tc.expectedResult) {
					t.Errorf("expected post %v, got %v", tc.expectedResult, res)
				}
			}
		})
	}
}
