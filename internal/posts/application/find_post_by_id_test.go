package application

import (
	"context"
	"errors"
	"testing"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
	mock_domain "github.com/jdgonzalez907/saas-api/mocks/posts/domain"
	"github.com/stretchr/testify/mock"
)

func TestFindPostByID_Execute(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repo *mock_domain.MockPostRepository)
		id       int64
		wantPost bool
		wantErr  error
	}{
		{
			name: "success",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "test-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
			},
			id:       1,
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "error - post not found",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				repo.On("FindByID", mock.Anything, int64(99)).Return(nil, nil)
			},
			id:       99,
			wantPost: false,
			wantErr:  domain.ErrFindPostByID,
		},
		{
			name: "error - repository error",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database connection failed"))
			},
			id:       1,
			wantPost: false,
			wantErr:  domain.ErrFindPostByID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockPostRepository(t)
			tt.setup(t, repo)
			uc := NewFindPostByID(repo)
			got, err := uc.Execute(context.Background(), tt.id)

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