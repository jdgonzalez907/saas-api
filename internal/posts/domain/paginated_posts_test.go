package domain

import (
	"testing"
	"time"
)

func TestPaginatedPostsVO(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	contentInfo, _ := NewContentInformation("Post Title", []Block{titleBlock})
	now := time.Now().UTC()

	post, err := NewPost(PostParams{
		ID:                 1,
		ContentInformation: contentInfo,
		Status:             StatusPublished,
		CreatedAt:          now,
		UpdatedAt:          now,
		AuthorID:           10,
		LastEditorID:       10,
		PublishedAt:        &now,
	})
	if err != nil {
		t.Fatalf("unexpected error creating post: %v", err)
	}

	posts := []*Post{post}
	nextPublishedAt := &now
	nextID := int64(1)

	vo := NewPaginatedPosts(posts, nextPublishedAt, &nextID)

	if len(vo.Posts()) != 1 {
		t.Errorf("expected 1 post, got %d", len(vo.Posts()))
	}
	if vo.Posts()[0].ID() != post.ID() {
		t.Errorf("expected post ID %d, got %d", post.ID(), vo.Posts()[0].ID())
	}
	if vo.NextPublishedAt() != nextPublishedAt {
		t.Errorf("expected NextPublishedAt %v, got %v", nextPublishedAt, vo.NextPublishedAt())
	}
	if vo.NextID() != &nextID {
		t.Errorf("expected NextID %v, got %v", &nextID, vo.NextID())
	}

	dto := vo.ToDTO()
	if len(dto.Posts) != 1 {
		t.Errorf("expected 1 post in DTO, got %d", len(dto.Posts))
	}
	if dto.Posts[0].ID != post.ID() {
		t.Errorf("expected DTO post ID %d, got %d", post.ID(), dto.Posts[0].ID)
	}
	if dto.NextPublishedAt != nextPublishedAt {
		t.Errorf("expected DTO NextPublishedAt %v, got %v", nextPublishedAt, dto.NextPublishedAt)
	}
	if dto.NextID != &nextID {
		t.Errorf("expected DTO NextID %v, got %v", &nextID, dto.NextID)
	}
}
