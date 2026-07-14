package domain

import (
	"errors"
	"testing"
)

func TestNewContentInformation_Validation(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Main title")
	paragraphBlock, _ := NewParagraphBlock("hello", nil)
	markBlock, _ := NewMarkBlock("bold", "world")
	listItemBlock, _ := NewListItemBlock("item", nil)

	testCases := []struct {
		name    string
		title   string
		content []Block
		wantErr error
	}{
		{
			name:    "fail - empty title",
			title:   "",
			content: []Block{titleBlock},
			wantErr: ErrEmptyPostTitle,
		},
		{
			name:    "fail - orphan mark block at root",
			title:   "Title",
			content: []Block{titleBlock, markBlock},
			wantErr: ErrOrphanBlock,
		},
		{
			name:    "fail - orphan listitem block at root",
			title:   "Title",
			content: []Block{listItemBlock},
			wantErr: ErrOrphanBlock,
		},
		{
			name:    "success - valid content",
			title:   "Title",
			content: []Block{titleBlock, paragraphBlock},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewContentInformation(tc.title, tc.content)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestContentInformation_Getters(t *testing.T) {
	titleBlock, _ := NewTitleBlock("Title")
	content := []Block{titleBlock}
	info, err := NewContentInformation("My Title", content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Title() != "My Title" {
		t.Errorf("expected title My Title, got %s", info.Title())
	}
	if len(info.Content()) != 1 {
		t.Errorf("expected content length 1, got %d", len(info.Content()))
	}
	if !info.Content()[0].Equals(titleBlock) {
		t.Error("expected first block to equal titleBlock")
	}
}

func TestContentInformation_Equals(t *testing.T) {
	titleA, _ := NewTitleBlock("Title A")
	titleB, _ := NewTitleBlock("Title B")

	info1, _ := NewContentInformation("My Title", []Block{titleA})
	info2, _ := NewContentInformation("My Title", []Block{titleA})
	infoDiffTitle, _ := NewContentInformation("Other Title", []Block{titleA})
	infoDiffContentLen, _ := NewContentInformation("My Title", nil)
	infoDiffContentVal, _ := NewContentInformation("My Title", []Block{titleB})

	testCases := []struct {
		name     string
		base     ContentInformation
		other    ContentInformation
		expected bool
	}{
		{
			name:     "success - identical content information",
			base:     info1,
			other:    info2,
			expected: true,
		},
		{
			name:     "fail - different title",
			base:     info1,
			other:    infoDiffTitle,
			expected: false,
		},
		{
			name:     "fail - different content length",
			base:     info1,
			other:    infoDiffContentLen,
			expected: false,
		},
		{
			name:     "fail - different content value",
			base:     info1,
			other:    infoDiffContentVal,
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

func TestContentInformationDTO_Mapping(t *testing.T) {
	titleBlock, _ := NewTitleBlock("My Title")
	info, _ := NewContentInformation("Header Title", []Block{titleBlock})

	dto := info.ToDTO()

	if dto.Title != "Header Title" {
		t.Errorf("expected dto title Header Title, got %s", dto.Title)
	}
	if len(dto.Content) != 1 {
		t.Fatalf("expected 1 content block DTO, got %d", len(dto.Content))
	}
	if dto.Content[0].Text != "My Title" {
		t.Errorf("expected content block text My Title, got %s", dto.Content[0].Text)
	}

	reconstructed, err := ContentInformationFromDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error reconstructing: %v", err)
	}

	if !info.Equals(reconstructed) {
		t.Error("expected reconstructed content information to equal original")
	}
}

func TestContentInformationFromDTO_Errors(t *testing.T) {
	dto := ContentInformationDTO{
		Title: "Header Title",
		Content: []BlockDTO{
			{
				Type: "invalid-type",
				Text: "text",
			},
		},
	}

	_, err := ContentInformationFromDTO(dto)
	if !errors.Is(err, ErrInvalidBlockType) {
		t.Errorf("expected ErrInvalidBlockType, got %v", err)
	}
}
