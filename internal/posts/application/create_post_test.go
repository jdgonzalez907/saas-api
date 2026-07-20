package application

import (
	"context"
	"errors"
	"testing"

	"github.com/jdgonzalez907/saas-api/internal/posts/domain"
	mock_domain "github.com/jdgonzalez907/saas-api/mocks/posts/domain"
	"github.com/stretchr/testify/mock"
)

func validPostInput() (string, string, string, []domain.Block, domain.PostStatus, int64) {
	return "Test Title", "test-slug", "http://example.com/cover.jpg", []domain.Block{}, domain.PostStatusDraft, 1
}

func mustNewPostWithSlug(t *testing.T, slug string) *domain.Post {
	t.Helper()
	blocks := []domain.Block{}
	p, err := domain.New("Test Title", slug, "http://example.com/cover.jpg", blocks, domain.PostStatusDraft, 1)
	if err != nil {
		t.Fatalf("mustNewPostWithSlug() error = %v", err)
	}
	return p
}

func mustNewAutor(t *testing.T, id int64) *domain.Autor {
	t.Helper()
	a, err := domain.NewAutor(id, "John Doe")
	if err != nil {
		t.Fatalf("mustNewAutor() error = %v", err)
	}
	return a
}

func TestCreatePost_Execute(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository)
		input   struct {
			title    string
			slug     string
			cover    string
			content  []domain.Block
			status   domain.PostStatus
			authorID int64
		}
		wantPost bool
		wantErr  error
	}{
		{
			name: "success - with blocks",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:   "Test Title",
				slug:    "test-slug",
				cover:   "http://example.com/cover.jpg",
				content: []domain.Block{},
				status:  domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "success - empty blocks",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  nil,
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "success - with paragraph block",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title: "Test Title",
				slug:  "test-slug",
				cover: "http://example.com/cover.jpg",
				content: func() []domain.Block {
					b, _ := domain.NewParagraphBlock("Test content", nil)
					return []domain.Block{b}
				}(),
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: true,
			wantErr:  nil,
		},
		{
			name: "error - autor not found",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrCreatePost,
		},
		{
			name: "error - FindByID autor fails",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrCreatePost,
		},
		{
			name: "error - slug already exists",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				existing := mustNewPostWithSlug(t, "test-slug")
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(existing, nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrCreatePost,
		},
		{
			name: "error - FindBySlug fails",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, errors.New("database error"))
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrCreatePost,
		},
		{
			name: "error - Create fails",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrCreatePost,
		},
		{
			name: "error - invalid title (empty)",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrPostTitleRequired,
		},
		{
			name: "error - invalid slug (empty)",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "").Return(nil, nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrPostSlugRequired,
		},
		{
			name: "error - invalid cover (empty)",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(1)).Return(mustNewAutor(t, 1), nil)
				repo.On("FindBySlug", mock.Anything, "test-slug").Return(nil, nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 1,
			},
			wantPost: false,
			wantErr:  domain.ErrPostCoverRequired,
		},
		{
			name: "error - invalid authorID",
			setup: func(t *testing.T, repo *mock_domain.MockPostRepository, autorRepo *mock_domain.MockAutorRepository) {
				autorRepo.On("FindByID", mock.Anything, int64(0)).Return(nil, nil)
			},
			input: struct {
				title    string
				slug     string
				cover    string
				content  []domain.Block
				status   domain.PostStatus
				authorID int64
			}{
				title:    "Test Title",
				slug:     "test-slug",
				cover:    "http://example.com/cover.jpg",
				content:  []domain.Block{},
				status:   domain.PostStatusDraft,
				authorID: 0,
			},
			wantPost: false,
			wantErr:  domain.ErrCreatePost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockPostRepository(t)
			autorRepo := mock_domain.NewMockAutorRepository(t)
			tt.setup(t, repo, autorRepo)
			uc := NewCreatePost(repo, autorRepo)
			got, err := uc.Execute(context.Background(), tt.input.title, tt.input.slug, tt.input.cover, tt.input.content, tt.input.status, tt.input.authorID)

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
