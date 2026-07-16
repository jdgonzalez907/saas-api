package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewPostStatus(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    PostStatus
		wantErr error
	}{
		{
			name:    "success - status draft",
			input:   "draft",
			want:    StatusDraft,
			wantErr: nil,
		},
		{
			name:    "success - status published",
			input:   "published",
			want:    StatusPublished,
			wantErr: nil,
		},
		{
			name:    "fail - invalid status",
			input:   "invalid",
			want:    "",
			wantErr: ErrInvalidPostStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewPostStatus(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
			if got != tc.want {
				t.Errorf("expected status %s, got %s", tc.want, got)
			}
		})
	}
}

func TestNewPost_Validation(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo, _ := NewContentInformation("Post Title", []Block{titleBlock})
	now := time.Now().UTC()

	testCases := []struct {
		name             string
		id               int64
		contentInfo      ContentInformation
		status           PostStatus
		createdAt        time.Time
		updatedAt        time.Time
		authorID         int64
		lastEditorID     int64
		publishedAt      *time.Time
		wantErr          error
	}{
		{
			name:         "fail - invalid ID (zero)",
			id:           0,
			contentInfo:  contentInfo,
			status:       StatusDraft,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: 1,
			wantErr:      ErrInvalidPostID,
		},
		{
			name:         "fail - invalid ID (negative)",
			id:           -5,
			contentInfo:  contentInfo,
			status:       StatusDraft,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: 1,
			wantErr:      ErrInvalidPostID,
		},
		{
			name:         "fail - invalid AuthorID",
			id:           1,
			contentInfo:  contentInfo,
			status:       StatusDraft,
			createdAt:    now,
			updatedAt:    now,
			authorID:     0,
			lastEditorID: 1,
			wantErr:      ErrInvalidAuthorID,
		},
		{
			name:         "fail - invalid LastEditorID",
			id:           1,
			contentInfo:  contentInfo,
			status:       StatusDraft,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: -1,
			wantErr:      ErrInvalidLastEditorID,
		},
		{
			name:         "fail - draft with publication date",
			id:           1,
			contentInfo:  contentInfo,
			status:       StatusDraft,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: 1,
			publishedAt:  &now,
			wantErr:      ErrDraftCannotHavePublicationDate,
		},
		{
			name:         "fail - published without publication date",
			id:           1,
			contentInfo:  contentInfo,
			status:       StatusPublished,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: 1,
			publishedAt:  nil,
			wantErr:      ErrPublishedMustHavePublicationDate,
		},
		{
			name:         "success - valid draft",
			id:           1,
			contentInfo:  contentInfo,
			status:       StatusDraft,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: 1,
			publishedAt:  nil,
			wantErr:      nil,
		},
		{
			name:         "success - valid published",
			id:           1,
			contentInfo:  contentInfo,
			status:       StatusPublished,
			createdAt:    now,
			updatedAt:    now,
			authorID:     1,
			lastEditorID: 1,
			publishedAt:  &now,
			wantErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPost(tc.id, tc.contentInfo, tc.status, tc.createdAt, tc.updatedAt, tc.authorID, tc.lastEditorID, tc.publishedAt)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestNewPostWithoutID_Validation(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo, _ := NewContentInformation("Post Title", []Block{titleBlock})

	t.Run("fail - invalid authorID", func(t *testing.T) {
		_, err := NewPostWithoutID(contentInfo, StatusDraft, 0)
		if !errors.Is(err, ErrInvalidAuthorID) {
			t.Errorf("expected ErrInvalidAuthorID, got %v", err)
		}
	})

	t.Run("success - valid draft", func(t *testing.T) {
		post, err := NewPostWithoutID(contentInfo, StatusDraft, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if post.ID() != UnassignedPostID {
			t.Errorf("expected ID to be unassigned (%d), got %d", UnassignedPostID, post.ID())
		}
		if post.AuthorID() != 5 {
			t.Errorf("expected AuthorID 5, got %d", post.AuthorID())
		}
		if post.LastEditorID() != 5 {
			t.Errorf("expected LastEditorID 5, got %d", post.LastEditorID())
		}
		if post.PublishedAt() != nil {
			t.Error("expected publishedAt to be nil for draft")
		}
		if post.CreatedAt().IsZero() || post.UpdatedAt().IsZero() {
			t.Error("expected timestamps to be initialized")
		}
	})

	t.Run("success - valid published", func(t *testing.T) {
		post, err := NewPostWithoutID(contentInfo, StatusPublished, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if post.PublishedAt() == nil {
			t.Fatal("expected publishedAt to be set for published post")
		}
		if post.PublishedAt().IsZero() {
			t.Error("expected publication date to be valid")
		}
	})
}

func TestPost_Getters(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title A")
	contentInfo, _ := NewContentInformation("Post A", []Block{titleBlock})
	now := time.Now().UTC()

	post, err := NewPost(1, contentInfo, StatusDraft, now, now, 10, 12, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID() != 1 {
		t.Errorf("expected ID 1, got %d", post.ID())
	}
	if !post.ContentInformation().Equals(contentInfo) {
		t.Error("expected content information to equal contentInfo")
	}
	if post.Status() != StatusDraft {
		t.Errorf("expected status %s, got %s", StatusDraft, post.Status())
	}
	if !post.CreatedAt().Equal(now) {
		t.Errorf("expected createdAt %v, got %v", now, post.CreatedAt())
	}
	if !post.UpdatedAt().Equal(now) {
		t.Errorf("expected updatedAt %v, got %v", now, post.UpdatedAt())
	}
	if post.AuthorID() != 10 {
		t.Errorf("expected AuthorID 10, got %d", post.AuthorID())
	}
	if post.LastEditorID() != 12 {
		t.Errorf("expected LastEditorID 12, got %d", post.LastEditorID())
	}
	if post.PublishedAt() != nil {
		t.Errorf("expected PublishedAt to be nil, got %v", post.PublishedAt())
	}
}

func TestPost_Equals(t *testing.T) {
	titleBlock1, _ := NewTitleBlock("Title A")
	titleBlock2, _ := NewTitleBlock("Title B")
	contentInfo1, _ := NewContentInformation("Post A", []Block{titleBlock1})
	contentInfo2, _ := NewContentInformation("Post A", []Block{titleBlock2})

	now := time.Now().UTC()
	anotherTime := now.Add(time.Second)

	post1, _ := NewPost(1, contentInfo1, StatusDraft, now, now, 10, 10, nil)

	post2, _ := NewPost(1, contentInfo1, StatusDraft, now, now, 10, 10, nil)

	postDifferentID, _ := NewPost(2, contentInfo1, StatusDraft, now, now, 10, 10, nil)

	postDifferentStatus, _ := NewPost(1, contentInfo1, StatusPublished, now, now, 10, 10, &now)

	postDifferentAuthorID, _ := NewPost(1, contentInfo1, StatusDraft, now, now, 11, 10, nil)

	postDifferentLastEditorID, _ := NewPost(1, contentInfo1, StatusDraft, now, now, 10, 11, nil)

	postDifferentCreatedAt, _ := NewPost(1, contentInfo1, StatusDraft, now.Add(time.Second), now, 10, 10, nil)

	postDifferentUpdatedAt, _ := NewPost(1, contentInfo1, StatusDraft, now, now.Add(time.Second), 10, 10, nil)

	postDifferentContentInfo, _ := NewPost(1, contentInfo2, StatusDraft, now, now, 10, 10, nil)

	postPublishedBase, _ := NewPost(1, contentInfo1, StatusPublished, now, now, 10, 10, &now)

	postPublishedDiffTime, _ := NewPost(1, contentInfo1, StatusPublished, now, now, 10, 10, &anotherTime)

	postPublishedDiffContent, _ := NewPost(1, contentInfo2, StatusPublished, now, now, 10, 10, &now)

	// Bypassing constructor to test direct Equals boundary checks for coverage
	postMismatchedPublishedAt1 := &Post{
		id:          1,
		status:      StatusPublished,
		publishedAt: nil,
	}
	postMismatchedPublishedAt2 := &Post{
		id:          1,
		status:      StatusPublished,
		publishedAt: &now,
	}

	testCases := []struct {
		name     string
		base     *Post
		other    *Post
		expected bool
	}{
		{
			name:     "success - identical draft posts",
			base:     post1,
			other:    post2,
			expected: true,
		},
		{
			name:     "fail - other is nil",
			base:     post1,
			other:    nil,
			expected: false,
		},
		{
			name:     "fail - different ID",
			base:     post1,
			other:    postDifferentID,
			expected: false,
		},
		{
			name:     "fail - different status",
			base:     post1,
			other:    postDifferentStatus,
			expected: false,
		},
		{
			name:     "fail - different authorID",
			base:     post1,
			other:    postDifferentAuthorID,
			expected: false,
		},
		{
			name:     "fail - different lastEditorID",
			base:     post1,
			other:    postDifferentLastEditorID,
			expected: false,
		},
		{
			name:     "fail - different createdAt",
			base:     post1,
			other:    postDifferentCreatedAt,
			expected: false,
		},
		{
			name:     "fail - different updatedAt",
			base:     post1,
			other:    postDifferentUpdatedAt,
			expected: false,
		},
		{
			name:     "fail - different contentInformation",
			base:     post1,
			other:    postDifferentContentInfo,
			expected: false,
		},
		{
			name:     "fail - draft vs published (nil vs non-nil publishedAt)",
			base:     post1,
			other:    postPublishedBase,
			expected: false,
		},
		{
			name:     "fail - different publishedAt times",
			base:     postPublishedBase,
			other:    postPublishedDiffTime,
			expected: false,
		},
		{
			name:     "success - identical published posts",
			base:     postPublishedBase,
			other:    postPublishedBase,
			expected: true,
		},
		{
			name:     "fail - published vs draft (non-nil vs nil publishedAt)",
			base:     postPublishedBase,
			other:    post1,
			expected: false,
		},
		{
			name:     "fail - published posts with different content",
			base:     postPublishedBase,
			other:    postPublishedDiffContent,
			expected: false,
		},
		{
			name:     "fail - same status mismatched publishedAt pointers",
			base:     postMismatchedPublishedAt1,
			other:    postMismatchedPublishedAt2,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.base.Equals(tc.other)
			if got != tc.expected {
				t.Errorf("expected Equals result to be %t, got %t", tc.expected, got)
			}
		})
	}
}

func TestPost_UpdateContentAndStatus(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo1, _ := NewContentInformation("Original", []Block{titleBlock})
	contentInfo2, _ := NewContentInformation("Updated", []Block{titleBlock})
	now := time.Now().UTC()

	post, _ := NewPost(1, contentInfo1, StatusDraft, now, now, 5, 5, nil)

	t.Run("fail - invalid lastEditorID", func(t *testing.T) {
		_, err := post.UpdateContentAndStatus(contentInfo2, StatusPublished, 0)
		if !errors.Is(err, ErrInvalidLastEditorID) {
			t.Errorf("expected ErrInvalidLastEditorID, got %v", err)
		}
	})

	t.Run("success - draft to published", func(t *testing.T) {
		time.Sleep(10 * time.Millisecond)
		updated, err := post.UpdateContentAndStatus(contentInfo2, StatusPublished, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if updated.ID() != post.ID() {
			t.Error("expected ID to remain unchanged")
		}
		if !updated.CreatedAt().Equal(post.CreatedAt()) {
			t.Error("expected CreatedAt to remain unchanged")
		}
		if updated.AuthorID() != post.AuthorID() {
			t.Error("expected AuthorID to remain unchanged")
		}

		if !updated.ContentInformation().Equals(contentInfo2) {
			t.Error("expected content to be updated")
		}
		if updated.Status() != StatusPublished {
			t.Errorf("expected status to be updated to %s, got %s", StatusPublished, updated.Status())
		}
		if updated.LastEditorID() != 6 {
			t.Errorf("expected LastEditorID to be 6, got %d", updated.LastEditorID())
		}
		if !updated.UpdatedAt().After(post.UpdatedAt()) {
			t.Error("expected UpdatedAt to be updated to a newer time")
		}
		if updated.PublishedAt() == nil {
			t.Error("expected PublishedAt to be set when transitioning to published")
		}
	})

	t.Run("success - published remains published (preserves publishedAt)", func(t *testing.T) {
		pubTime := now.Add(-time.Hour)
		pubPost, _ := NewPost(1, contentInfo1, StatusPublished, now, now, 5, 5, &pubTime)

		updated, err := pubPost.UpdateContentAndStatus(contentInfo2, StatusPublished, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !updated.PublishedAt().Equal(pubTime) {
			t.Errorf("expected publishedAt to be preserved as %v, got %v", pubTime, updated.PublishedAt())
		}
	})

	t.Run("success - published to draft (clears publishedAt)", func(t *testing.T) {
		pubPost, _ := NewPost(1, contentInfo1, StatusPublished, now, now, 5, 5, &now)

		updated, err := pubPost.UpdateContentAndStatus(contentInfo2, StatusDraft, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if updated.Status() != StatusDraft {
			t.Errorf("expected status to be draft, got %s", updated.Status())
		}
		if updated.PublishedAt() != nil {
			t.Error("expected publishedAt to be cleared (nil) when transitioning back to draft")
		}
	})
}

func TestPostDTO_Mapping(t *testing.T) {
	titleBlock, _ := NewTitleBlock("My Title")
	contentInfo, _ := NewContentInformation("My Header", []Block{titleBlock})
	now := time.Now().UTC()

	post, _ := NewPost(15, contentInfo, StatusPublished, now, now, 8, 9, &now)

	dto := post.ToDTO()

	if dto.ID != 15 {
		t.Errorf("expected dto ID 15, got %d", dto.ID)
	}
	if dto.Title != "My Header" {
		t.Errorf("expected dto Title My Header, got %s", dto.Title)
	}
	if dto.Status != "published" {
		t.Errorf("expected dto Status published, got %s", dto.Status)
	}
	if dto.AuthorID != 8 {
		t.Errorf("expected dto AuthorID 8, got %d", dto.AuthorID)
	}
	if dto.LastEditorID != 9 {
		t.Errorf("expected dto LastEditorID 9, got %d", dto.LastEditorID)
	}
	if dto.PublishedAt == nil || !dto.PublishedAt.Equal(now) {
		t.Errorf("expected dto PublishedAt %v, got %v", now, dto.PublishedAt)
	}

	// Normal reconstruction
	reconstructed, err := PostFromDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error reconstructing: %v", err)
	}

	if !post.Equals(reconstructed) {
		t.Error("expected reconstructed post to equal original")
	}

	// Reconstruction with UnassignedPostID (0)
	dto.ID = UnassignedPostID
	reconstructedUnassigned, err := PostFromDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error reconstructing unassigned post: %v", err)
	}
	if reconstructedUnassigned.ID() != UnassignedPostID {
		t.Errorf("expected reconstructed unassigned post ID %d, got %d", UnassignedPostID, reconstructedUnassigned.ID())
	}
}

func TestPostFromDTO_Validation(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name        string
		dto         *PostDTO
		wantErr     error
		expectNil   bool
		expectUnass bool
	}{
		{
			name:      "success - nil dto returns nil",
			dto:       nil,
			wantErr:   nil,
			expectNil: true,
		},
		{
			name: "fail - invalid ContentInformation DTO",
			dto: &PostDTO{
				ID: 1,
				ContentInformationDTO: ContentInformationDTO{
					Title: "Title",
					Content: []BlockDTO{
						{
							Type: "invalid-type",
						},
					},
				},
				Status:       "draft",
				CreatedAt:    now,
				UpdatedAt:    now,
				AuthorID:     1,
				LastEditorID: 1,
			},
			wantErr: ErrInvalidBlockType,
		},
		{
			name: "fail - invalid PostStatus DTO",
			dto: &PostDTO{
				ID: 1,
				ContentInformationDTO: ContentInformationDTO{
					Title:   "Title",
					Content: nil,
				},
				Status:       "invalid-status",
				CreatedAt:    now,
				UpdatedAt:    now,
				AuthorID:     1,
				LastEditorID: 1,
			},
			wantErr: ErrInvalidPostStatus,
		},
		{
			name: "fail - unassigned post with draft state having published date",
			dto: &PostDTO{
				ID: UnassignedPostID,
				ContentInformationDTO: ContentInformationDTO{
					Title:   "Title",
					Content: nil,
				},
				Status:       "draft",
				CreatedAt:    now,
				UpdatedAt:    now,
				AuthorID:     1,
				LastEditorID: 1,
				PublishedAt:  &now,
			},
			wantErr: ErrDraftCannotHavePublicationDate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PostFromDTO(tc.dto)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
			if tc.expectNil && got != nil {
				t.Error("expected returned post to be nil")
			}
		})
	}
}

func TestPost_EnsureInvariants_Private(t *testing.T) {
	t.Run("ensureInvariants - negative ID", func(t *testing.T) {
		p := &Post{
			id: -1,
		}
		err := p.ensureInvariants()
		if !errors.Is(err, ErrInvalidPostID) {
			t.Errorf("expected ErrInvalidPostID, got %v", err)
		}
	})
}

func TestPost_AssignID(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo, _ := NewContentInformation("Post Title", []Block{titleBlock})
	post, _ := NewPostWithoutID(contentInfo, StatusDraft, 10)

	if post.ID() != UnassignedPostID {
		t.Errorf("expected ID to be unassigned (0), got %d", post.ID())
	}

	post.AssignID(42)
	if post.ID() != 42 {
		t.Errorf("expected ID to be 42 after AssignID, got %d", post.ID())
	}
}
