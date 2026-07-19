package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
	mock_domain "github.com/jdgonzalez907/saas-api/mocks/posts/domain"
	"github.com/stretchr/testify/mock"
)

func mustNewPostForUpdate(t *testing.T, id, authorID int64, slug string) *domain.Post {
	t.Helper()
	now := time.Now()
	titleBlock, _ := domain.NewTitleBlock("Test Title", nil)
	content := []domain.Block{titleBlock}
	p, err := domain.NewWithID(id, "Test Title", slug, "http://example.com/cover.jpg", content, domain.PostStatusDraft, authorID, now, now)
	if err != nil {
		t.Fatalf("mustNewPostForUpdate() error = %v", err)
	}
	return p
}

func TestUpdateContent_Execute(t *testing.T) {
	content := func() []domain.Block {
		b, _ := domain.NewParagraphBlock("New content", nil)
		return []domain.Block{b}
	}()

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mock_domain.MockPostRepository)
		id      int64
		input   struct {
			executedBy int64
			title      string
			slug       string
			cover      string
			content    []domain.Block
			status     domain.PostStatus
		}
		wantPost bool
		wantErr  error
	}{
		{
			name: "success - same slug",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "test-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "test-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "success - slug changed",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "old-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				repo.On("FindBySlug", mock.Anything, "new-slug").Return(nil, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "new-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "error - FindByID fails",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "test-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrUpdateContent,
		},
		{
			name: "error - post not found",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "test-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrUpdateContent,
		},
		{
			name: "error - slug already exists",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "old-slug")
				existing := mustNewPostForUpdate(t, 2, 99, "new-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				repo.On("FindBySlug", mock.Anything, "new-slug").Return(existing, nil)
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "new-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrUpdateContent,
		},
		{
			name: "error - FindBySlug fails",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "old-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				repo.On("FindBySlug", mock.Anything, "new-slug").Return(nil, errors.New("database error"))
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "new-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrUpdateContent,
		},
		{
			name: "error - unauthorized",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "test-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 99,
				title:      "Updated Title",
				slug:       "test-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrPostUnauthorizedUpdate,
		},
		{
			name: "error - invalid title",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "test-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "",
				slug:       "test-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrPostTitleRequired,
		},
		{
			name: "error - Update fails",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository) {
				post := mustNewPostForUpdate(t, 1, 1, "test-slug")
				repo.On("FindByID", mock.Anything, int64(1)).Return(post, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			id: 1,
			input: struct {
				executedBy int64
				title      string
				slug       string
				cover      string
				content    []domain.Block
				status     domain.PostStatus
			}{
				executedBy: 1,
				title:      "Updated Title",
				slug:       "test-slug",
				cover:      "http://example.com/new-cover.jpg",
				content:    content,
				status:     domain.PostStatusPublished,
			},
			wantPost: false,
			wantErr:  domain.ErrUpdateContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockPostRepository(t)
			tt.setup(t, repo)
			uc := NewUpdateContent(repo)
			got, err := uc.Execute(context.Background(), tt.id, tt.input.executedBy, tt.input.title, tt.input.slug, tt.input.cover, tt.input.content, tt.input.status)

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
