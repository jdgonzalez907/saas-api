package domain

import (
	"errors"
	"strings"
)

type BlockType string

const (
	TypeTitle     BlockType = "title"
	TypeSubtitle  BlockType = "subtitle"
	TypeParagraph BlockType = "paragraph"
	TypeList      BlockType = "list"
	TypeListItem  BlockType = "listitem"
	TypeImage     BlockType = "image"
	TypeMark      BlockType = "mark"
)

var (
	ErrInvalidBlockType          = errors.New("invalid content block type")
	ErrInvalidBlockStructure     = errors.New("block structure does not match its type")
	ErrEmptyImageURL             = errors.New("image block must specify a valid image url")
	ErrInvalidMarkStyle          = errors.New("text formatting style must be bold or italic")
	ErrInvalidMarkStructure      = errors.New("formatted text cannot contain nested blocks")
	ErrInvalidListitemStructure  = errors.New("list item can only contain formatted text")
	ErrInvalidParagraphStructure = errors.New("paragraph can only contain formatted text")
	ErrInvalidTitleStructure     = errors.New("title block cannot contain nested blocks")
	ErrInvalidSubtitleStructure  = errors.New("subtitle block cannot contain nested blocks")
	ErrListDirectTextNotEmpty    = errors.New("list block cannot contain direct text, only list items")
	ErrInvalidListChildren       = errors.New("list block can only contain list items")
)

type Block struct {
	blockType BlockType
	text      string
	children  []Block
}

type BlockDTO struct {
	Type     string     `json:"type"`
	Text     string     `json:"text"`
	Children []BlockDTO `json:"children"`
}

func NewBlock(blockType BlockType, text string, children []Block) (Block, error) {
	switch blockType {
	case TypeTitle, TypeSubtitle, TypeParagraph, TypeList, TypeListItem, TypeImage, TypeMark:
		// Valid type
	default:
		return Block{}, ErrInvalidBlockType
	}

	switch blockType {
	case TypeTitle:
		if len(children) > 0 {
			return Block{}, ErrInvalidTitleStructure
		}
	case TypeSubtitle:
		if len(children) > 0 {
			return Block{}, ErrInvalidSubtitleStructure
		}
	case TypeImage:
		if len(children) > 0 {
			return Block{}, ErrInvalidBlockStructure
		}
		if text == "" {
			return Block{}, ErrEmptyImageURL
		}
	case TypeListItem:
		for _, child := range children {
			if child.blockType != TypeMark {
				return Block{}, ErrInvalidListitemStructure
			}
		}
	case TypeParagraph:
		for _, child := range children {
			if child.blockType != TypeMark {
				return Block{}, ErrInvalidParagraphStructure
			}
		}
	case TypeList:
		if text != "" {
			return Block{}, ErrListDirectTextNotEmpty
		}
		for _, child := range children {
			if child.blockType != TypeListItem {
				return Block{}, ErrInvalidListChildren
			}
		}
	case TypeMark:
		if len(children) > 0 {
			return Block{}, ErrInvalidMarkStructure
		}
		if !strings.HasPrefix(text, "bold:") && !strings.HasPrefix(text, "italic:") {
			return Block{}, ErrInvalidMarkStyle
		}
	}

	return Block{
		blockType: blockType,
		text:      text,
		children:  children,
	}, nil
}

func NewTitleBlock(text string) (Block, error) {
	return NewBlock(TypeTitle, text, nil)
}

func NewSubtitleBlock(text string) (Block, error) {
	return NewBlock(TypeSubtitle, text, nil)
}

func NewParagraphBlock(text string, marks []Block) (Block, error) {
	return NewBlock(TypeParagraph, text, marks)
}

func NewListBlock(items []Block) (Block, error) {
	return NewBlock(TypeList, "", items)
}

func NewListItemBlock(text string, marks []Block) (Block, error) {
	return NewBlock(TypeListItem, text, marks)
}

func NewImageBlock(url string) (Block, error) {
	return NewBlock(TypeImage, url, nil)
}

func NewMarkBlock(style string, content string) (Block, error) {
	if style != "bold" && style != "italic" {
		return Block{}, ErrInvalidMarkStyle
	}
	return NewBlock(TypeMark, style+":"+content, nil)
}

func (b Block) Type() BlockType {
	return b.blockType
}

func (b Block) Text() string {
	return b.text
}

func (b Block) Children() []Block {
	return b.children
}

func (b Block) Equals(other Block) bool {
	if b.blockType != other.blockType {
		return false
	}
	if b.text != other.text {
		return false
	}
	if len(b.children) != len(other.children) {
		return false
	}
	for i := range b.children {
		if !b.children[i].Equals(other.children[i]) {
			return false
		}
	}
	return true
}

func (b Block) ToDTO() BlockDTO {
	var childrenDTO []BlockDTO
	if b.children != nil {
		childrenDTO = make([]BlockDTO, len(b.children))
		for i, child := range b.children {
			childrenDTO[i] = child.ToDTO()
		}
	}
	return BlockDTO{
		Type:     string(b.blockType),
		Text:     b.text,
		Children: childrenDTO,
	}
}

func BlockFromDTO(dto BlockDTO) (Block, error) {
	var children []Block
	if dto.Children != nil {
		children = make([]Block, len(dto.Children))
		for i, childDTO := range dto.Children {
			child, err := BlockFromDTO(childDTO)
			if err != nil {
				return Block{}, err
			}
			children[i] = child
		}
	}
	return NewBlock(BlockType(dto.Type), dto.Text, children)
}
