package domain

import "testing"

func mustNewBlock(t *testing.T, blockType BlockType, value string, children []Block) Block {
	t.Helper()
	b, err := NewBlock(blockType, value, children)
	if err != nil {
		t.Fatalf("NewBlock() error = %v", err)
	}
	return b
}

func mustNewBold(t *testing.T, value string) Block {
	t.Helper()
	b, err := NewBoldBlock(value)
	if err != nil {
		t.Fatalf("NewBoldBlock() error = %v", err)
	}
	return b
}

func mustNewItalic(t *testing.T, value string) Block {
	t.Helper()
	b, err := NewItalicBlock(value)
	if err != nil {
		t.Fatalf("NewItalicBlock() error = %v", err)
	}
	return b
}

func mustNewLinkURL(t *testing.T, value string) Block {
	t.Helper()
	b, err := NewLinkURLBlock(value)
	if err != nil {
		t.Fatalf("NewLinkURLBlock() error = %v", err)
	}
	return b
}

func mustNewLink(t *testing.T, value string, child Block) Block {
	t.Helper()
	b, err := NewLinkBlock(value, []Block{child})
	if err != nil {
		t.Fatalf("NewLinkBlock() error = %v", err)
	}
	return b
}

func mustNewListItem(t *testing.T, value string, children []Block) Block {
	t.Helper()
	b, err := NewListItemBlock(value, children)
	if err != nil {
		t.Fatalf("NewListItemBlock() error = %v", err)
	}
	return b
}

func TestNewBlock(t *testing.T) {
	tests := []struct {
		name      string
		blockType BlockType
		value     string
		children  []Block
		wantErr   error
	}{
		{
			name:      "success - title",
			blockType: BlockTypeTitle,
			value:     "Hello",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - subtitle",
			blockType: BlockTypeSubtitle,
			value:     "World",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - image",
			blockType: BlockTypeImage,
			value:     "img.png",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - paragraph",
			blockType: BlockTypeParagraph,
			value:     "Text",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - list",
			blockType: BlockTypeList,
			value:     "",
			children:  []Block{mustNewListItem(t, "item1", nil)},
			wantErr:   nil,
		},
		{
			name:      "success - list_item",
			blockType: BlockTypeListItem,
			value:     "item",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - bold",
			blockType: BlockTypeBold,
			value:     "bold",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - italic",
			blockType: BlockTypeItalic,
			value:     "italic",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "success - link",
			blockType: BlockTypeLink,
			value:     "click",
			children:  []Block{mustNewLinkURL(t, "https://example.com")},
			wantErr:   nil,
		},
		{
			name:      "success - link_url",
			blockType: BlockTypeLinkURL,
			value:     "https://example.com",
			children:  nil,
			wantErr:   nil,
		},
		{
			name:      "error - invalid type",
			blockType: "invalid",
			value:     "test",
			children:  nil,
			wantErr:   ErrBlockInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBlock(tt.blockType, tt.value, tt.children)
			if err != tt.wantErr {
				t.Errorf("NewBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != tt.blockType {
					t.Errorf("NewBlock() BlockType = %v, want %v", b.BlockType(), tt.blockType)
				}
				if b.Value() != tt.value {
					t.Errorf("NewBlock() Value = %v, want %v", b.Value(), tt.value)
				}
			}
		})
	}
}

func TestNewTitleBlock(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		children []Block
		wantErr  error
	}{
		{
			name:     "success - without children",
			value:    "Title",
			children: nil,
			wantErr:  nil,
		},
		{
			name:     "success - with italic child",
			value:    "Title with",
			children: []Block{mustNewItalic(t, "italic")},
			wantErr:  nil,
		},
		{
			name:     "error - empty value",
			value:    "",
			children: nil,
			wantErr:  ErrBlockEmptyValue,
		},
		{
			name:     "error - invalid child type bold",
			value:    "Title",
			children: []Block{mustNewBold(t, "bold")},
			wantErr:  ErrBlockInvalidChildren,
		},
		{
			name:     "error - invalid child type link",
			value:    "Title",
			children: []Block{mustNewLink(t, "link", mustNewLinkURL(t, "url"))},
			wantErr:  ErrBlockInvalidChildren,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewTitleBlock(tt.value, tt.children)
			if err != tt.wantErr {
				t.Errorf("NewTitleBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeTitle {
					t.Errorf("NewTitleBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeTitle)
				}
				if b.Value() != tt.value {
					t.Errorf("NewTitleBlock() Value = %v, want %v", b.Value(), tt.value)
				}
			}
		})
	}
}

func TestNewSubtitleBlock(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		children []Block
		wantErr  error
	}{
		{
			name:     "success - without children",
			value:    "Subtitle",
			children: nil,
			wantErr:  nil,
		},
		{
			name:     "success - with italic child",
			value:    "Subtitle",
			children: []Block{mustNewItalic(t, "italic")},
			wantErr:  nil,
		},
		{
			name:     "success - with bold child",
			value:    "Subtitle",
			children: []Block{mustNewBold(t, "bold")},
			wantErr:  nil,
		},
		{
			name: "success - with italic and bold children",
			value: "Subtitle",
			children: []Block{
				mustNewItalic(t, "italic"),
				mustNewBold(t, "bold"),
			},
			wantErr: nil,
		},
		{
			name:     "error - empty value",
			value:    "",
			children: nil,
			wantErr:  ErrBlockEmptyValue,
		},
		{
			name:     "error - invalid child type link",
			value:    "Subtitle",
			children: []Block{mustNewLink(t, "link", mustNewLinkURL(t, "url"))},
			wantErr:  ErrBlockInvalidChildren,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewSubtitleBlock(tt.value, tt.children)
			if err != tt.wantErr {
				t.Errorf("NewSubtitleBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeSubtitle {
					t.Errorf("NewSubtitleBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeSubtitle)
				}
			}
		})
	}
}

func TestNewImageBlock(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "success",
			value:   "image.png",
			wantErr: nil,
		},
		{
			name:    "error - empty value",
			value:   "",
			wantErr: ErrBlockEmptyValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewImageBlock(tt.value)
			if err != tt.wantErr {
				t.Errorf("NewImageBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeImage {
					t.Errorf("NewImageBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeImage)
				}
				if b.Value() != tt.value {
					t.Errorf("NewImageBlock() Value = %v, want %v", b.Value(), tt.value)
				}
			}
		})
	}
}

func TestNewParagraphBlock(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		children []Block
		wantErr  error
	}{
		{
			name:     "success - without children",
			value:    "Paragraph",
			children: nil,
			wantErr:  nil,
		},
		{
			name:     "success - with bold child",
			value:    "Paragraph",
			children: []Block{mustNewBold(t, "bold")},
			wantErr:  nil,
		},
		{
			name:     "success - with italic child",
			value:    "Paragraph",
			children: []Block{mustNewItalic(t, "italic")},
			wantErr:  nil,
		},
		{
			name:     "success - with link child",
			value:    "Paragraph",
			children: []Block{mustNewLink(t, "link", mustNewLinkURL(t, "url"))},
			wantErr:  nil,
		},
		{
			name: "success - with multiple valid children",
			value: "Paragraph",
			children: []Block{
				mustNewBold(t, "bold"),
				mustNewItalic(t, "italic"),
				mustNewLink(t, "link", mustNewLinkURL(t, "url")),
			},
			wantErr: nil,
		},
		{
			name:     "error - empty value",
			value:    "",
			children: nil,
			wantErr:  ErrBlockEmptyValue,
		},
		{
			name:     "error - invalid child type title",
			value:    "Paragraph",
			children: []Block{mustNewBlock(t, BlockTypeTitle, "title", nil)},
			wantErr:  ErrBlockInvalidChildren,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewParagraphBlock(tt.value, tt.children)
			if err != tt.wantErr {
				t.Errorf("NewParagraphBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeParagraph {
					t.Errorf("NewParagraphBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeParagraph)
				}
			}
		})
	}
}

func TestNewListBlock(t *testing.T) {
	tests := []struct {
		name     string
		children []Block
		wantErr  error
	}{
		{
			name: "success - with one item",
			children: []Block{
				mustNewListItem(t, "item1", nil),
			},
			wantErr: nil,
		},
		{
			name: "success - with multiple items",
			children: []Block{
				mustNewListItem(t, "item1", nil),
				mustNewListItem(t, "item2", nil),
			},
			wantErr: nil,
		},
		{
			name:     "error - empty children",
			children: []Block{},
			wantErr:  ErrBlockListRequiresItems,
		},
		{
			name:     "error - nil children",
			children: nil,
			wantErr:  ErrBlockListRequiresItems,
		},
		{
			name: "error - invalid child type",
			children: []Block{
				mustNewBold(t, "bold"),
			},
			wantErr: ErrBlockInvalidChildren,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewListBlock(tt.children)
			if err != tt.wantErr {
				t.Errorf("NewListBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeList {
					t.Errorf("NewListBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeList)
				}
				if len(b.Children()) != len(tt.children) {
					t.Errorf("NewListBlock() Children length = %v, want %v", len(b.Children()), len(tt.children))
				}
			}
		})
	}
}

func TestNewListItemBlock(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		children []Block
		wantErr  error
	}{
		{
			name:     "success - without children",
			value:    "item",
			children: nil,
			wantErr:  nil,
		},
		{
			name:     "success - with italic child",
			value:    "item",
			children: []Block{mustNewItalic(t, "italic")},
			wantErr:  nil,
		},
		{
			name:     "success - with bold child",
			value:    "item",
			children: []Block{mustNewBold(t, "bold")},
			wantErr:  nil,
		},
		{
			name:     "success - with link child",
			value:    "item",
			children: []Block{mustNewLink(t, "link", mustNewLinkURL(t, "url"))},
			wantErr:  nil,
		},
		{
			name: "success - with all valid children",
			value: "item",
			children: []Block{
				mustNewBold(t, "bold"),
				mustNewItalic(t, "italic"),
				mustNewLink(t, "link", mustNewLinkURL(t, "url")),
			},
			wantErr: nil,
		},
		{
			name:     "error - empty value",
			value:    "",
			children: nil,
			wantErr:  ErrBlockEmptyValue,
		},
		{
			name:     "error - invalid child type title",
			value:    "item",
			children: []Block{mustNewBlock(t, BlockTypeTitle, "title", nil)},
			wantErr:  ErrBlockInvalidChildren,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewListItemBlock(tt.value, tt.children)
			if err != tt.wantErr {
				t.Errorf("NewListItemBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeListItem {
					t.Errorf("NewListItemBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeListItem)
				}
			}
		})
	}
}

func TestNewBoldBlock(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "success",
			value:   "bold text",
			wantErr: nil,
		},
		{
			name:    "error - empty value",
			value:   "",
			wantErr: ErrBlockEmptyValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBoldBlock(tt.value)
			if err != tt.wantErr {
				t.Errorf("NewBoldBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeBold {
					t.Errorf("NewBoldBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeBold)
				}
				if b.Value() != tt.value {
					t.Errorf("NewBoldBlock() Value = %v, want %v", b.Value(), tt.value)
				}
			}
		})
	}
}

func TestNewItalicBlock(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "success",
			value:   "italic text",
			wantErr: nil,
		},
		{
			name:    "error - empty value",
			value:   "",
			wantErr: ErrBlockEmptyValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewItalicBlock(tt.value)
			if err != tt.wantErr {
				t.Errorf("NewItalicBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeItalic {
					t.Errorf("NewItalicBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeItalic)
				}
				if b.Value() != tt.value {
					t.Errorf("NewItalicBlock() Value = %v, want %v", b.Value(), tt.value)
				}
			}
		})
	}
}

func TestNewLinkBlock(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		children []Block
		wantErr  error
	}{
		{
			name:     "success",
			value:    "click here",
			children: []Block{mustNewLinkURL(t, "https://example.com")},
			wantErr:  nil,
		},
		{
			name:     "error - empty value",
			value:    "",
			children: []Block{mustNewLinkURL(t, "https://example.com")},
			wantErr:  ErrBlockEmptyValue,
		},
		{
			name:     "error - no children",
			value:    "click",
			children: nil,
			wantErr:  ErrBlockLinkRequiresChild,
		},
		{
			name:     "error - empty children",
			value:    "click",
			children: []Block{},
			wantErr:  ErrBlockLinkRequiresChild,
		},
		{
			name:     "error - more than one child",
			value:    "click",
			children: []Block{mustNewLinkURL(t, "url1"), mustNewLinkURL(t, "url2")},
			wantErr:  ErrBlockLinkRequiresChild,
		},
		{
			name:     "error - child is not link_url",
			value:    "click",
			children: []Block{mustNewBold(t, "bold")},
			wantErr:  ErrBlockLinkRequiresChild,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewLinkBlock(tt.value, tt.children)
			if err != tt.wantErr {
				t.Errorf("NewLinkBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeLink {
					t.Errorf("NewLinkBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeLink)
				}
				if b.Value() != tt.value {
					t.Errorf("NewLinkBlock() Value = %v, want %v", b.Value(), tt.value)
				}
				if len(b.Children()) != 1 {
					t.Errorf("NewLinkBlock() Children length = %v, want 1", len(b.Children()))
				}
			}
		})
	}
}

func TestNewLinkURLBlock(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{
			name:    "success",
			value:   "https://example.com",
			wantErr: nil,
		},
		{
			name:    "error - empty value",
			value:   "",
			wantErr: ErrBlockEmptyValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewLinkURLBlock(tt.value)
			if err != tt.wantErr {
				t.Errorf("NewLinkURLBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if b.BlockType() != BlockTypeLinkURL {
					t.Errorf("NewLinkURLBlock() BlockType = %v, want %v", b.BlockType(), BlockTypeLinkURL)
				}
				if b.Value() != tt.value {
					t.Errorf("NewLinkURLBlock() Value = %v, want %v", b.Value(), tt.value)
				}
			}
		})
	}
}

func TestBlock_Equals(t *testing.T) {
	tests := []struct {
		name string
		v1   Block
		v2   Block
		want bool
	}{
		{
			name: "equal - same title blocks",
			v1:   mustNewBlock(t, BlockTypeTitle, "Title", nil),
			v2:   mustNewBlock(t, BlockTypeTitle, "Title", nil),
			want: true,
		},
		{
			name: "not equal - different type",
			v1:   mustNewBlock(t, BlockTypeTitle, "Text", nil),
			v2:   mustNewBlock(t, BlockTypeSubtitle, "Text", nil),
			want: false,
		},
		{
			name: "not equal - different value",
			v1:   mustNewBlock(t, BlockTypeTitle, "Title1", nil),
			v2:   mustNewBlock(t, BlockTypeTitle, "Title2", nil),
			want: false,
		},
		{
			name: "not equal - different children length",
			v1:   mustNewBlock(t, BlockTypeList, "", []Block{mustNewListItem(t, "item", nil)}),
			v2: mustNewBlock(t, BlockTypeList, "", []Block{
				mustNewListItem(t, "item1", nil),
				mustNewListItem(t, "item2", nil),
			}),
			want: false,
		},
		{
			name: "not equal - different children values",
			v1:   mustNewBlock(t, BlockTypeList, "", []Block{mustNewListItem(t, "item1", nil)}),
			v2:   mustNewBlock(t, BlockTypeList, "", []Block{mustNewListItem(t, "item2", nil)}),
			want: false,
		},
		{
			name: "equal - nested blocks",
			v1: mustNewBlock(t, BlockTypeList, "", []Block{
				mustNewListItem(t, "item", []Block{mustNewBold(t, "bold")}),
			}),
			v2: mustNewBlock(t, BlockTypeList, "", []Block{
				mustNewListItem(t, "item", []Block{mustNewBold(t, "bold")}),
			}),
			want: true,
		},
		{
			name: "equal - both empty",
			v1:   Block{},
			v2:   Block{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Equals(tt.v2); got != tt.want {
				t.Errorf("Block.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_ToDTO(t *testing.T) {
	tests := []struct {
		name string
		block Block
		want BlockDTO
	}{
		{
			name:  "title without children",
			block: mustNewBlock(t, BlockTypeTitle, "Title", nil),
			want: BlockDTO{
				BlockType: "title",
				Value:     "Title",
				Children:  nil,
			},
		},
		{
			name:  "bold block",
			block: mustNewBold(t, "bold"),
			want: BlockDTO{
				BlockType: "bold",
				Value:     "bold",
				Children:  nil,
			},
		},
		{
			name: "link with child",
			block: mustNewLink(t, "click", mustNewLinkURL(t, "https://example.com")),
			want: BlockDTO{
				BlockType: "link",
				Value:     "click",
				Children: []BlockDTO{
					{
						BlockType: "link_url",
						Value:     "https://example.com",
						Children:  nil,
					},
				},
			},
		},
		{
			name: "list with items",
			block: mustNewBlock(t, BlockTypeList, "", []Block{
				mustNewListItem(t, "item1", nil),
				mustNewListItem(t, "item2", nil),
			}),
			want: BlockDTO{
				BlockType: "list",
				Value:     "",
				Children: []BlockDTO{
					{BlockType: "list_item", Value: "item1", Children: nil},
					{BlockType: "list_item", Value: "item2", Children: nil},
				},
			},
		},
		{
			name:  "empty block",
			block: Block{},
			want:  BlockDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.block.ToDTO()
			if got.BlockType != tt.want.BlockType {
				t.Errorf("ToDTO().BlockType = %v, want %v", got.BlockType, tt.want.BlockType)
			}
			if got.Value != tt.want.Value {
				t.Errorf("ToDTO().Value = %v, want %v", got.Value, tt.want.Value)
			}
			if len(got.Children) != len(tt.want.Children) {
				t.Errorf("ToDTO().Children length = %v, want %v", len(got.Children), len(tt.want.Children))
				return
			}
			for i := range got.Children {
				if got.Children[i].BlockType != tt.want.Children[i].BlockType {
					t.Errorf("ToDTO().Children[%d].BlockType = %v, want %v", i, got.Children[i].BlockType, tt.want.Children[i].BlockType)
				}
				if got.Children[i].Value != tt.want.Children[i].Value {
					t.Errorf("ToDTO().Children[%d].Value = %v, want %v", i, got.Children[i].Value, tt.want.Children[i].Value)
				}
			}
		})
	}
}
