package domain

import (
	"testing"
	"time"
)

func validPostContent(t *testing.T) []Block {
	t.Helper()
	titleBlock, err := NewTitleBlock("Title", nil)
	if err != nil {
		t.Fatalf("NewTitleBlock() error = %v", err)
	}
	return []Block{titleBlock}
}

func TestNewPost(t *testing.T) {
	content := validPostContent(t)

	tests := []struct {
		name      string
		title     string
		slug      string
		cover     string
		content   []Block
		status    PostStatus
		authorID  int64
		wantErr   error
	}{
		{
			name:     "success",
			title:    "My Post",
			slug:     "my-post",
			cover:    "cover.png",
			content:  content,
			status:   PostStatusDraft,
			authorID: 1,
			wantErr:  nil,
		},
		{
			name:     "success - with nil content",
			title:    "My Post",
			slug:     "my-post",
			cover:    "cover.png",
			content:  nil,
			status:   PostStatusDraft,
			authorID: 1,
			wantErr:  nil,
		},
		{
			name:     "error - empty title",
			title:    "",
			slug:     "my-post",
			cover:    "cover.png",
			content:  content,
			status:   PostStatusDraft,
			authorID: 1,
			wantErr:  ErrPostTitleRequired,
		},
		{
			name:     "error - empty slug",
			title:    "My Post",
			slug:     "",
			cover:    "cover.png",
			content:  content,
			status:   PostStatusDraft,
			authorID: 1,
			wantErr:  ErrPostSlugRequired,
		},
		{
			name:     "error - empty cover",
			title:    "My Post",
			slug:     "my-post",
			cover:    "",
			content:  content,
			status:   PostStatusDraft,
			authorID: 1,
			wantErr:  ErrPostCoverRequired,
		},
		{
			name:     "error - empty status",
			title:    "My Post",
			slug:     "my-post",
			cover:    "cover.png",
			content:  content,
			status:   "",
			authorID: 1,
			wantErr:  ErrPostStatusRequired,
		},
		{
			name:     "error - zero authorID",
			title:    "My Post",
			slug:     "my-post",
			cover:    "cover.png",
			content:  content,
			status:   PostStatusDraft,
			authorID: 0,
			wantErr:  ErrPostAuthorIDRequired,
		},
		{
			name:     "error - negative authorID",
			title:    "My Post",
			slug:     "my-post",
			cover:    "cover.png",
			content:  content,
			status:   PostStatusDraft,
			authorID: -1,
			wantErr:  ErrPostAuthorIDRequired,
		},
		{
			name:    "error - invalid root block",
			title:   "My Post",
			slug:    "my-post",
			cover:   "cover.png",
			content: []Block{mustNewBold(t, "bold")},
			status:  PostStatusDraft,
			authorID: 1,
			wantErr: ErrPostInvalidRootBlock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.title, tt.slug, tt.cover, tt.content, tt.status, tt.authorID)
			if err != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if p.ID() != 0 {
					t.Errorf("New().ID() = %v, want %v", p.ID(), 0)
				}
				if p.Title() != tt.title {
					t.Errorf("New().Title() = %v, want %v", p.Title(), tt.title)
				}
				if p.Slug() != tt.slug {
					t.Errorf("New().Slug() = %v, want %v", p.Slug(), tt.slug)
				}
				if p.Cover() != tt.cover {
					t.Errorf("New().Cover() = %v, want %v", p.Cover(), tt.cover)
				}
				if p.Status() != tt.status {
					t.Errorf("New().Status() = %v, want %v", p.Status(), tt.status)
				}
				if p.AuthorID() != tt.authorID {
					t.Errorf("New().AuthorID() = %v, want %v", p.AuthorID(), tt.authorID)
				}
				if p.CreatedAt().IsZero() {
					t.Errorf("New().CreatedAt() should not be zero")
				}
				if p.UpdatedAt().IsZero() {
					t.Errorf("New().UpdatedAt() should not be zero")
				}
			}
		})
	}
}

func TestNewWithID(t *testing.T) {
	content := validPostContent(t)
	now := time.Now()

	tests := []struct {
		name      string
		id        int64
		title     string
		slug      string
		cover     string
		content   []Block
		status    PostStatus
		authorID  int64
		createdAt time.Time
		updatedAt time.Time
		wantErr   error
	}{
		{
			name:      "success",
			id:        1,
			title:     "My Post",
			slug:      "my-post",
			cover:     "cover.png",
			content:   content,
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   nil,
		},
		{
			name:      "error - zero ID",
			id:        0,
			title:     "My Post",
			slug:      "my-post",
			cover:     "cover.png",
			content:   content,
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostIDRequired,
		},
		{
			name:      "error - negative ID",
			id:        -1,
			title:     "My Post",
			slug:      "my-post",
			cover:     "cover.png",
			content:   content,
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostIDRequired,
		},
		{
			name:      "error - empty title",
			id:        1,
			title:     "",
			slug:      "my-post",
			cover:     "cover.png",
			content:   content,
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostTitleRequired,
		},
		{
			name:      "error - empty slug",
			id:        1,
			title:     "My Post",
			slug:      "",
			cover:     "cover.png",
			content:   content,
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostSlugRequired,
		},
		{
			name:      "error - empty cover",
			id:        1,
			title:     "My Post",
			slug:      "my-post",
			cover:     "",
			content:   content,
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostCoverRequired,
		},
		{
			name:      "error - empty status",
			id:        1,
			title:     "My Post",
			slug:      "my-post",
			cover:     "cover.png",
			content:   content,
			status:    "",
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostStatusRequired,
		},
		{
			name:      "error - zero authorID",
			id:        1,
			title:     "My Post",
			slug:      "my-post",
			cover:     "cover.png",
			content:   content,
			status:    PostStatusDraft,
			authorID:  0,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostAuthorIDRequired,
		},
		{
			name:      "error - invalid root block",
			id:        1,
			title:     "My Post",
			slug:      "my-post",
			cover:     "cover.png",
			content:   []Block{mustNewBold(t, "bold")},
			status:    PostStatusDraft,
			authorID:  1,
			createdAt: now,
			updatedAt: now,
			wantErr:   ErrPostInvalidRootBlock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewWithID(tt.id, tt.title, tt.slug, tt.cover, tt.content, tt.status, tt.authorID, tt.createdAt, tt.updatedAt)
			if err != tt.wantErr {
				t.Errorf("NewWithID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if p.ID() != tt.id {
					t.Errorf("NewWithID().ID() = %v, want %v", p.ID(), tt.id)
				}
				if p.Title() != tt.title {
					t.Errorf("NewWithID().Title() = %v, want %v", p.Title(), tt.title)
				}
				if p.Slug() != tt.slug {
					t.Errorf("NewWithID().Slug() = %v, want %v", p.Slug(), tt.slug)
				}
				if p.Cover() != tt.cover {
					t.Errorf("NewWithID().Cover() = %v, want %v", p.Cover(), tt.cover)
				}
				if p.Status() != tt.status {
					t.Errorf("NewWithID().Status() = %v, want %v", p.Status(), tt.status)
				}
				if p.AuthorID() != tt.authorID {
					t.Errorf("NewWithID().AuthorID() = %v, want %v", p.AuthorID(), tt.authorID)
				}
				if !p.CreatedAt().Equal(tt.createdAt) {
					t.Errorf("NewWithID().CreatedAt() = %v, want %v", p.CreatedAt(), tt.createdAt)
				}
				if !p.UpdatedAt().Equal(tt.updatedAt) {
					t.Errorf("NewWithID().UpdatedAt() = %v, want %v", p.UpdatedAt(), tt.updatedAt)
				}
			}
		})
	}
}

func TestPost_AssignID(t *testing.T) {
	p := mustNewPost(t, "My Post", "my-post", "cover.png", validPostContent(t), PostStatusDraft, 1)

	p.AssignID(42)

	if got := p.ID(); got != 42 {
		t.Errorf("Post.ID() = %v, want %v", got, 42)
	}
}

func TestPost_Equals(t *testing.T) {
	content := validPostContent(t)
	now := time.Now()

	tests := []struct {
		name string
		e1   *Post
		e2   *Post
		want bool
	}{
		{
			name: "equal - same ID",
			e1:   mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			e2:   mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			want: true,
		},
		{
			name: "not equal - different ID",
			e1:   mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			e2:   mustNewPostWithID(t, 2, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			want: false,
		},
		{
			name: "not equal - nil other",
			e1:   mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			e2:   nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e1.Equals(tt.e2); got != tt.want {
				t.Errorf("Post.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost_UpdateContent(t *testing.T) {
	content := validPostContent(t)
	newContent := []Block{mustNewBlock(t, BlockTypeParagraph, "New content", nil)}
	now := time.Now()

	tests := []struct {
		name       string
		entity     *Post
		title      string
		slug       string
		cover      string
		content    []Block
		status     PostStatus
		executedBy int64
		wantErr    error
	}{
		{
			name:       "success",
			entity:     mustNewPostWithID(t, 1, "Old Title", "old-slug", "old-cover", content, PostStatusDraft, 1, now, now),
			title:      "New Title",
			slug:       "new-slug",
			cover:      "new-cover.png",
			content:    newContent,
			status:     PostStatusPublished,
			executedBy: 1,
			wantErr:    nil,
		},
		{
			name:       "error - unauthorized",
			entity:     mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			title:      "New Title",
			slug:       "new-slug",
			cover:      "new-cover.png",
			content:    newContent,
			status:     PostStatusPublished,
			executedBy: 2,
			wantErr:    ErrPostUnauthorizedUpdate,
		},
		{
			name:       "error - empty title",
			entity:     mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			title:      "",
			slug:       "new-slug",
			cover:      "new-cover.png",
			content:    newContent,
			status:     PostStatusPublished,
			executedBy: 1,
			wantErr:    ErrPostTitleRequired,
		},
		{
			name:       "error - empty slug",
			entity:     mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			title:      "New Title",
			slug:       "",
			cover:      "new-cover.png",
			content:    newContent,
			status:     PostStatusPublished,
			executedBy: 1,
			wantErr:    ErrPostSlugRequired,
		},
		{
			name:       "error - empty cover",
			entity:     mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			title:      "New Title",
			slug:       "new-slug",
			cover:      "",
			content:    newContent,
			status:     PostStatusPublished,
			executedBy: 1,
			wantErr:    ErrPostCoverRequired,
		},
		{
			name:       "error - empty status",
			entity:     mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			title:      "New Title",
			slug:       "new-slug",
			cover:      "new-cover.png",
			content:    newContent,
			status:     "",
			executedBy: 1,
			wantErr:    ErrPostStatusRequired,
		},
		{
			name:       "error - invalid root block",
			entity:     mustNewPostWithID(t, 1, "Title", "slug", "cover", content, PostStatusDraft, 1, now, now),
			title:      "New Title",
			slug:       "new-slug",
			cover:      "new-cover.png",
			content:    []Block{mustNewBold(t, "bold")},
			status:     PostStatusPublished,
			executedBy: 1,
			wantErr:    ErrPostInvalidRootBlock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := tt.entity.UpdatedAt()
			err := tt.entity.UpdateContent(tt.title, tt.slug, tt.cover, tt.content, tt.status, tt.executedBy)
			if err != tt.wantErr {
				t.Errorf("Post.UpdateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got := tt.entity.Title(); got != tt.title {
					t.Errorf("Post.Title() = %v, want %v", got, tt.title)
				}
				if got := tt.entity.Slug(); got != tt.slug {
					t.Errorf("Post.Slug() = %v, want %v", got, tt.slug)
				}
				if got := tt.entity.Cover(); got != tt.cover {
					t.Errorf("Post.Cover() = %v, want %v", got, tt.cover)
				}
				if got := tt.entity.Status(); got != tt.status {
					t.Errorf("Post.Status() = %v, want %v", got, tt.status)
				}
				if tt.entity.UpdatedAt().Before(updatedAt) {
					t.Errorf("Post.UpdatedAt() should be updated")
				}
			}
		})
	}
}

func TestPost_ToDTO(t *testing.T) {
	content := validPostContent(t)
	now := time.Now()
	p := mustNewPostWithID(t, 1, "My Post", "my-post", "cover.png", content, PostStatusDraft, 1, now, now)

	dto := p.ToDTO()

	if dto.ID != 1 {
		t.Errorf("Post.ToDTO().ID = %v, want %v", dto.ID, 1)
	}
	if dto.Title != "My Post" {
		t.Errorf("Post.ToDTO().Title = %v, want %v", dto.Title, "My Post")
	}
	if dto.Slug != "my-post" {
		t.Errorf("Post.ToDTO().Slug = %v, want %v", dto.Slug, "my-post")
	}
	if dto.Cover != "cover.png" {
		t.Errorf("Post.ToDTO().Cover = %v, want %v", dto.Cover, "cover.png")
	}
	if dto.Status != string(PostStatusDraft) {
		t.Errorf("Post.ToDTO().Status = %v, want %v", dto.Status, PostStatusDraft)
	}
	if dto.AuthorID != 1 {
		t.Errorf("Post.ToDTO().AuthorID = %v, want %v", dto.AuthorID, 1)
	}
	if !dto.CreatedAt.Equal(now) {
		t.Errorf("Post.ToDTO().CreatedAt = %v, want %v", dto.CreatedAt, now)
	}
	if !dto.UpdatedAt.Equal(now) {
		t.Errorf("Post.ToDTO().UpdatedAt = %v, want %v", dto.UpdatedAt, now)
	}
	if len(dto.Content) != 1 {
		t.Errorf("Post.ToDTO().Content length = %v, want %v", len(dto.Content), 1)
	}
}

func TestPost_ToDTO_NilContent(t *testing.T) {
	now := time.Now()
	p := mustNewPostWithID(t, 1, "My Post", "my-post", "cover.png", nil, PostStatusDraft, 1, now, now)

	dto := p.ToDTO()

	if dto.Content != nil {
		t.Errorf("Post.ToDTO().Content should be nil, got %v", dto.Content)
	}
}

func mustNewPost(t *testing.T, title, slug, cover string, content []Block, status PostStatus, authorID int64) *Post {
	t.Helper()
	p, err := New(title, slug, cover, content, status, authorID)
	if err != nil {
		t.Fatalf("mustNewPost() error = %v", err)
	}
	return p
}

func mustNewPostWithID(t *testing.T, id int64, title, slug, cover string, content []Block, status PostStatus, authorID int64, createdAt, updatedAt time.Time) *Post {
	t.Helper()
	p, err := NewWithID(id, title, slug, cover, content, status, authorID, createdAt, updatedAt)
	if err != nil {
		t.Fatalf("mustNewPostWithID() error = %v", err)
	}
	return p
}
