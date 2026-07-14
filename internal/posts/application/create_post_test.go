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

func TestCreatePostUseCase(t *testing.T) {
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
		input            *domain.Post
		mockExpectations func(*domainMocks.MockPostRepository)
		expectedError    error
	}{
		{
			name:  "success - create post",
			input: post,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("Create", mock.Anything, post).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "fail - repository create error",
			input: post,
			mockExpectations: func(m *domainMocks.MockPostRepository) {
				m.On("Create", mock.Anything, post).Return(dbErr)
			},
			expectedError: domain.ErrCreatingPost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPostRepository := new(domainMocks.MockPostRepository)
			tc.mockExpectations(mockPostRepository)

			useCase := application.NewCreatePostUseCase(mockPostRepository)
			err := useCase.Execute(context.Background(), tc.input)

			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error: %v, got nil", tc.expectedError)
				}
				if !errors.Is(err, tc.expectedError) {
					if !(tc.expectedError == domain.ErrCreatingPost && errors.Unwrap(err) != nil) {
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
