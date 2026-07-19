package application

import (
	"context"
	"errors"
	"testing"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
	mock_domain "github.com/jdgonzalez907/saas-api/mocks/posts/domain"
	"github.com/stretchr/testify/mock"
)

func TestFindPostBySlug_Execute(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repo *mock_domain.MockPostRepository)
		slug     string
		wantPost bool
		wantErr  error
	}{
		{
			name: "success",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostWithSlug(t, "test-slug")
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(post, nil)
			},
			slug:     "test-slug",
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "error - post not found",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				repo.On("FindBySlug", mock.Anything, "nonexistent-slug").Return(nil, nil)
			},
			slug:     "nonexistent-slug",
			wantPost: false,
			wantErr:  domain.ErrFindPostBySlug,
		},
		{
			name: "error - repository error",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				repo.On("FindBySlug", mock.Anything, "error-slug").Return(nil, errors.New("database connection failed"))
			},
			slug:     "error-slug",
			wantPost: false,
			wantErr:  domain.ErrFindPostBySlug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockPostRepository(t)
			tt.setup(t, repo)
			uc := NewFindPostBySlug(repo)
			got, err := uc.Execute(context.Background(), tt.slug)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantPost && got == nil {
				t.Errorf("Execute() returned nil post, want non-nil")
			}

			if !tt.wantPost && got != nil {
				t.Errorf("Execute() returned non-nil post, want nil")
			}
		})
	}
}