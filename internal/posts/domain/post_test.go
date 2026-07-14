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
		name    string
		params  PostParams
		wantErr error
	}{
		{
			name: "fail - invalid ID",
			params: PostParams{
				ID:                 0,
				ContentInformation: contentInfo,
				Status:             StatusDraft,
				CreatedAt:          now,
				UpdatedAt:          now,
				CreatedBy:          1,
				UpdatedBy:          1,
			},
			wantErr: ErrInvalidPostID,
		},
		{
			name: "fail - invalid CreatedBy",
			params: PostParams{
				ID:                 1,
				ContentInformation: contentInfo,
				Status:             StatusDraft,
				CreatedAt:          now,
				UpdatedAt:          now,
				CreatedBy:          0,
				UpdatedBy:          1,
			},
			wantErr: ErrInvalidUserID,
		},
		{
			name: "fail - invalid UpdatedBy",
			params: PostParams{
				ID:                 1,
				ContentInformation: contentInfo,
				Status:             StatusDraft,
				CreatedAt:          now,
				UpdatedAt:          now,
				CreatedBy:          1,
				UpdatedBy:          -1,
			},
			wantErr: ErrInvalidUserID,
		},
		{
			name: "success - valid params",
			params: PostParams{
				ID:                 1,
				ContentInformation: contentInfo,
				Status:             StatusDraft,
				CreatedAt:          now,
				UpdatedAt:          now,
				CreatedBy:          1,
				UpdatedBy:          1,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPost(tc.params)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestNewPostWithoutID_Validation(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo, _ := NewContentInformation("Post Title", []Block{titleBlock})

	t.Run("fail - invalid createdBy", func(t *testing.T) {
		_, err := NewPostWithoutID(contentInfo, StatusDraft, 0)
		if !errors.Is(err, ErrInvalidUserID) {
			t.Errorf("expected ErrInvalidUserID, got %v", err)
		}
	})

	t.Run("success - valid", func(t *testing.T) {
		post, err := NewPostWithoutID(contentInfo, StatusDraft, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if post.ID() != UnassignedPostID {
			t.Errorf("expected ID to be unassigned (%d), got %d", UnassignedPostID, post.ID())
		}
		if post.CreatedBy() != 5 {
			t.Errorf("expected CreatedBy 5, got %d", post.CreatedBy())
		}
		if post.UpdatedBy() != 5 {
			t.Errorf("expected UpdatedBy 5, got %d", post.UpdatedBy())
		}
		if post.CreatedAt().IsZero() || post.UpdatedAt().IsZero() {
			t.Error("expected timestamps to be initialized")
		}
	})
}

func TestPost_Getters(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title A")
	contentInfo, _ := NewContentInformation("Post A", []Block{titleBlock})
	now := time.Now().UTC()

	post, err := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})
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
	if post.CreatedBy() != 10 {
		t.Errorf("expected CreatedBy 10, got %d", post.CreatedBy())
	}
	if post.UpdatedBy() != 10 {
		t.Errorf("expected UpdatedBy 10, got %d", post.UpdatedBy())
	}
}

func TestPost_Equals(t *testing.T) {
	titleBlock1, _ := NewTitleBlock("Title A")
	titleBlock2, _ := NewTitleBlock("Title B")
	contentInfo1, _ := NewContentInformation("Post A", []Block{titleBlock1})
	contentInfo2, _ := NewContentInformation("Post A", []Block{titleBlock2})

	now := time.Now().UTC()

	post1, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	post2, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	postDifferentID, _ := NewPost(PostParams{
		ID:                 2,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	postDifferentStatus, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusPublished,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	postDifferentCreatedBy, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          11,
		UpdatedBy:          10,
	})

	postDifferentUpdatedBy, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          11,
	})

	postDifferentCreatedAt, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now.Add(time.Second),
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	postDifferentUpdatedAt, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now.Add(time.Second),
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	postDifferentContentInfo, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo2,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          10,
		UpdatedBy:          10,
	})

	testCases := []struct {
		name     string
		base     *Post
		other    *Post
		expected bool
	}{
		{
			name:     "success - identical posts",
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
			name:     "fail - different createdBy",
			base:     post1,
			other:    postDifferentCreatedBy,
			expected: false,
		},
		{
			name:     "fail - different updatedBy",
			base:     post1,
			other:    postDifferentUpdatedBy,
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

func TestPost_Wither(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo1, _ := NewContentInformation("Original", []Block{titleBlock})
	contentInfo2, _ := NewContentInformation("Updated", []Block{titleBlock})
	now := time.Now().UTC()

	post, _ := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo1,
		Status:             StatusDraft,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          5,
		UpdatedBy:          5,
	})

	t.Run("fail - invalid updatedBy user ID", func(t *testing.T) {
		_, err := post.WithContentAndStatus(contentInfo2, StatusPublished, 0)
		if !errors.Is(err, ErrInvalidUserID) {
			t.Errorf("expected ErrInvalidUserID, got %v", err)
		}
	})

	t.Run("success - valid wither", func(t *testing.T) {
		time.Sleep(10 * time.Millisecond) // Ensure time moves forward
		updated, err := post.WithContentAndStatus(contentInfo2, StatusPublished, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if updated.ID() != post.ID() {
			t.Error("expected ID to remain unchanged")
		}
		if !updated.CreatedAt().Equal(post.CreatedAt()) {
			t.Error("expected CreatedAt to remain unchanged")
		}
		if updated.CreatedBy() != post.CreatedBy() {
			t.Error("expected CreatedBy to remain unchanged")
		}

		if !updated.ContentInformation().Equals(contentInfo2) {
			t.Error("expected content to be updated")
		}
		if updated.Status() != StatusPublished {
			t.Errorf("expected status to be updated to %s, got %s", StatusPublished, updated.Status())
		}
		if updated.UpdatedBy() != 6 {
			t.Errorf("expected UpdatedBy to be 6, got %d", updated.UpdatedBy())
		}
		if !updated.UpdatedAt().After(post.UpdatedAt()) {
			t.Error("expected UpdatedAt to be updated to a newer time")
		}
	})
}

func TestPostDTO_Mapping(t *testing.T) {
	titleBlock, _ := NewTitleBlock("My Title")
	contentInfo, _ := NewContentInformation("My Header", []Block{titleBlock})
	now := time.Now().UTC()

	post, _ := NewPost(PostParams{
		ID:                 15,
		ContentInformation: contentInfo,
		Status:             StatusPublished,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          8,
		UpdatedBy:          9,
	})

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
	if dto.CreatedBy != 8 {
		t.Errorf("expected dto CreatedBy 8, got %d", dto.CreatedBy)
	}
	if dto.UpdatedBy != 9 {
		t.Errorf("expected dto UpdatedBy 9, got %d", dto.UpdatedBy)
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
				Status:    "draft",
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: 1,
				UpdatedBy: 1,
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
				Status:    "invalid-status",
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: 1,
				UpdatedBy: 1,
			},
			wantErr: ErrInvalidPostStatus,
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
