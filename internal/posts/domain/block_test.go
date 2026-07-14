package domain

import (
	"errors"
	"testing"
)

func TestNewBlock_Validation(t *testing.T) {
	testCases := []struct {
		name      string
		blockType BlockType
		text      string
		children  []Block
		wantErr   error
	}{
		{
			name:      "fail - invalid block type",
			blockType: BlockType("invalid"),
			text:      "some text",
			children:  nil,
			wantErr:   ErrInvalidBlockType,
		},
		{
			name:      "success - valid title block",
			blockType: TypeTitle,
			text:      "Main Title",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "fail - title block with children",
			blockType: TypeTitle,
			text:      "Main Title",
			children:  []Block{{blockType: TypeMark, text: "bold:hello", children: nil}},
			wantErr:   ErrInvalidTitleStructure,
		},
		{
			name:      "success - valid subtitle block",
			blockType: TypeSubtitle,
			text:      "Subtitle",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "fail - subtitle block with children",
			blockType: TypeSubtitle,
			text:      "Subtitle",
			children:  []Block{{blockType: TypeMark, text: "bold:hello", children: nil}},
			wantErr:   ErrInvalidSubtitleStructure,
		},
		{
			name:      "success - valid image block",
			blockType: TypeImage,
			text:      "https://example.com/img.png",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "fail - image block with children",
			blockType: TypeImage,
			text:      "https://example.com/img.png",
			children:  []Block{{blockType: TypeMark, text: "bold:hello", children: nil}},
			wantErr:   ErrInvalidBlockStructure,
		},
		{
			name:      "fail - image block with empty url",
			blockType: TypeImage,
			text:      "",
			children:  nil,
			wantErr:   ErrEmptyImageURL,
		},
		{
			name:      "success - valid listitem block with mark children",
			blockType: TypeListItem,
			text:      "list item text",
			children:  []Block{{blockType: TypeMark, text: "bold:hello", children: nil}},
			wantErr:   nil,
		},
		{
			name:      "fail - listitem block with non-mark children",
			blockType: TypeListItem,
			text:      "list item text",
			children:  []Block{{blockType: TypeTitle, text: "title", children: nil}},
			wantErr:   ErrInvalidListitemStructure,
		},
		{
			name:      "success - valid paragraph block with mark children",
			blockType: TypeParagraph,
			text:      "paragraph text",
			children:  []Block{{blockType: TypeMark, text: "italic:hello", children: nil}},
			wantErr:   nil,
		},
		{
			name:      "fail - paragraph block with non-mark children",
			blockType: TypeParagraph,
			text:      "paragraph text",
			children:  []Block{{blockType: TypeSubtitle, text: "sub", children: nil}},
			wantErr:   ErrInvalidParagraphStructure,
		},
		{
			name:      "success - valid list block with listitem children",
			blockType: TypeList,
			text:      "",
			children:  []Block{{blockType: TypeListItem, text: "item", children: nil}},
			wantErr:   nil,
		},
		{
			name:      "fail - list block with direct text",
			blockType: TypeList,
			text:      "direct text",
			children:  nil,
			wantErr:   ErrListDirectTextNotEmpty,
		},
		{
			name:      "fail - list block with non-listitem children",
			blockType: TypeList,
			text:      "",
			children:  []Block{{blockType: TypeParagraph, text: "item", children: nil}},
			wantErr:   ErrInvalidListChildren,
		},
		{
			name:      "success - valid mark bold block",
			blockType: TypeMark,
			text:      "bold:bold text",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - valid mark italic block",
			blockType: TypeMark,
			text:      "italic:italic text",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "fail - mark block with children",
			blockType: TypeMark,
			text:      "bold:text",
			children:  []Block{{blockType: TypeMark, text: "italic:nested", children: nil}},
			wantErr:   ErrInvalidMarkStructure,
		},
		{
			name:      "fail - mark block with invalid format",
			blockType: TypeMark,
			text:      "underlined:text",
			children:  nil,
			wantErr:   ErrInvalidMarkStyle,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewBlock(tc.blockType, tc.text, tc.children)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestSpecificConstructors(t *testing.T) {
	t.Run("success - Title constructor", func(t *testing.T) {
		block, err := NewTitleBlock("Hello Title")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeTitle {
			t.Errorf("expected type %v, got %v", TypeTitle, block.Type())
		}
		if block.Text() != "Hello Title" {
			t.Errorf("expected text %s, got %s", "Hello Title", block.Text())
		}
		if len(block.Children()) != 0 {
			t.Errorf("expected empty children, got len %d", len(block.Children()))
		}
	})

	t.Run("success - Subtitle constructor", func(t *testing.T) {
		block, err := NewSubtitleBlock("Hello Subtitle")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeSubtitle {
			t.Errorf("expected type %v, got %v", TypeSubtitle, block.Type())
		}
		if block.Text() != "Hello Subtitle" {
			t.Errorf("expected text %s, got %s", "Hello Subtitle", block.Text())
		}
	})

	t.Run("success - Paragraph constructor", func(t *testing.T) {
		mark, _ := NewMarkBlock("bold", "marked")
		block, err := NewParagraphBlock("text", []Block{mark})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeParagraph {
			t.Errorf("expected type %v, got %v", TypeParagraph, block.Type())
		}
		if len(block.Children()) != 1 {
			t.Errorf("expected 1 child, got %d", len(block.Children()))
		}
	})

	t.Run("success - List constructor", func(t *testing.T) {
		item, _ := NewListItemBlock("item text", nil)
		block, err := NewListBlock([]Block{item})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeList {
			t.Errorf("expected type %v, got %v", TypeList, block.Type())
		}
	})

	t.Run("success - ListItem constructor", func(t *testing.T) {
		mark, _ := NewMarkBlock("italic", "nested")
		block, err := NewListItemBlock("text", []Block{mark})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeListItem {
			t.Errorf("expected type %v, got %v", TypeListItem, block.Type())
		}
	})

	t.Run("success - Image constructor", func(t *testing.T) {
		block, err := NewImageBlock("https://url.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeImage {
			t.Errorf("expected type %v, got %v", TypeImage, block.Type())
		}
	})

	t.Run("success - Mark constructor", func(t *testing.T) {
		block, err := NewMarkBlock("bold", "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if block.Type() != TypeMark {
			t.Errorf("expected type %v, got %v", TypeMark, block.Type())
		}
		if block.Text() != "bold:hello" {
			t.Errorf("expected formatted text bold:hello, got %s", block.Text())
		}
	})

	t.Run("fail - Mark constructor invalid style", func(t *testing.T) {
		_, err := NewMarkBlock("underlined", "hello")
		if !errors.Is(err, ErrInvalidMarkStyle) {
			t.Errorf("expected ErrInvalidMarkStyle, got %v", err)
		}
	})
}

func TestBlock_Equals(t *testing.T) {
	titleA, _ := NewTitleBlock("Title A")
	titleB, _ := NewTitleBlock("Title B")
	subtitleA, _ := NewSubtitleBlock("Title A")

	mark1, _ := NewMarkBlock("bold", "hello")
	mark2, _ := NewMarkBlock("italic", "hello")

	para1, _ := NewParagraphBlock("text", []Block{mark1})
	para2, _ := NewParagraphBlock("text", []Block{mark1})
	para3, _ := NewParagraphBlock("text", []Block{mark2})
	para4, _ := NewParagraphBlock("text", nil)

	testCases := []struct {
		name     string
		base     Block
		other    Block
		expected bool
	}{
		{
			name:     "success - identical title blocks",
			base:     titleA,
			other:    titleA,
			expected: true,
		},
		{
			name:     "fail - different type",
			base:     titleA,
			other:    subtitleA,
			expected: false,
		},
		{
			name:     "fail - different text",
			base:     titleA,
			other:    titleB,
			expected: false,
		},
		{
			name:     "success - identical paragraph blocks with children",
			base:     para1,
			other:    para2,
			expected: true,
		},
		{
			name:     "fail - different child content",
			base:     para1,
			other:    para3,
			expected: false,
		},
		{
			name:     "fail - different children count",
			base:     para1,
			other:    para4,
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

func TestBlockDTO_Mapping(t *testing.T) {
	mark, _ := NewMarkBlock("bold", "hello")
	para, _ := NewParagraphBlock("para text", []Block{mark})

	dto := para.ToDTO()

	if dto.Type != string(TypeParagraph) {
		t.Errorf("expected dto type %s, got %s", TypeParagraph, dto.Type)
	}
	if dto.Text != "para text" {
		t.Errorf("expected dto text para text, got %s", dto.Text)
	}
	if len(dto.Children) != 1 {
		t.Fatalf("expected 1 child dto, got %d", len(dto.Children))
	}
	if dto.Children[0].Type != string(TypeMark) {
		t.Errorf("expected child dto type %s, got %s", TypeMark, dto.Children[0].Type)
	}
	if dto.Children[0].Text != "bold:hello" {
		t.Errorf("expected child dto text bold:hello, got %s", dto.Children[0].Text)
	}

	reconstructed, err := BlockFromDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error reconstructing: %v", err)
	}

	if !para.Equals(reconstructed) {
		t.Error("expected reconstructed block to equal original block")
	}
}

func TestBlockFromDTO_Errors(t *testing.T) {
	dto := BlockDTO{
		Type: "paragraph",
		Text: "hello",
		Children: []BlockDTO{
			{
				Type: "invalid-child-type",
				Text: "child",
			},
		},
	}

	_, err := BlockFromDTO(dto)
	if !errors.Is(err, ErrInvalidBlockType) {
		t.Errorf("expected ErrInvalidBlockType, got %v", err)
	}
}
