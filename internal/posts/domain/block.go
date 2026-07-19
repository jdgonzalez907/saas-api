package domain

import (
	"errors"
)

const (
	linkRequiredChildrenCount = 1
)

var (
	ErrBlockInvalidType       = errors.New("block type is invalid")
	ErrBlockEmptyValue        = errors.New("block value is required")
	ErrBlockInvalidChildren   = errors.New("block has invalid children for its type")
	ErrBlockLinkRequiresChild = errors.New("link block requires exactly one link_url child")
	ErrBlockListRequiresItems = errors.New("list block requires at least one list_item child")
)

type BlockType string

const (
	BlockTypeTitle     BlockType = "title"
	BlockTypeSubtitle  BlockType = "subtitle"
	BlockTypeImage     BlockType = "image"
	BlockTypeParagraph BlockType = "paragraph"
	BlockTypeList      BlockType = "list"
	BlockTypeListItem  BlockType = "list_item"
	BlockTypeBold      BlockType = "bold"
	BlockTypeItalic    BlockType = "italic"
	BlockTypeLink      BlockType = "link"
	BlockTypeLinkURL   BlockType = "link_url"
)

type Block struct {
	blockType BlockType
	value     string
	children  []Block
}

type BlockDTO struct {
	BlockType string     `json:"block_type"`
	Value     string     `json:"value"`
	Children  []BlockDTO `json:"children,omitempty"`
}

func NewBlock(blockType BlockType, value string, children []Block) (Block, error) {
	switch blockType {
	case BlockTypeTitle:
		return NewTitleBlock(value, children)
	case BlockTypeSubtitle:
		return NewSubtitleBlock(value, children)
	case BlockTypeImage:
		return NewImageBlock(value)
	case BlockTypeParagraph:
		return NewParagraphBlock(value, children)
	case BlockTypeList:
		return NewListBlock(children)
	case BlockTypeListItem:
		return NewListItemBlock(value, children)
	case BlockTypeBold:
		return NewBoldBlock(value)
	case BlockTypeItalic:
		return NewItalicBlock(value)
	case BlockTypeLink:
		return NewLinkBlock(value, children)
	case BlockTypeLinkURL:
		return NewLinkURLBlock(value)
	default:
		return Block{}, ErrBlockInvalidType
	}
}

func NewTitleBlock(value string, children []Block) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	if err := validateChildren(children, []BlockType{BlockTypeItalic}); err != nil {
		return Block{}, err
	}

	return Block{blockType: BlockTypeTitle, value: value, children: children}, nil
}

func NewSubtitleBlock(value string, children []Block) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	if err := validateChildren(children, []BlockType{BlockTypeItalic, BlockTypeBold}); err != nil {
		return Block{}, err
	}

	return Block{blockType: BlockTypeSubtitle, value: value, children: children}, nil
}

func NewImageBlock(value string) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	return Block{blockType: BlockTypeImage, value: value}, nil
}

func NewParagraphBlock(value string, children []Block) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	if err := validateChildren(children, []BlockType{BlockTypeLink, BlockTypeBold, BlockTypeItalic}); err != nil {
		return Block{}, err
	}

	return Block{blockType: BlockTypeParagraph, value: value, children: children}, nil
}

func NewListBlock(children []Block) (Block, error) {
	if len(children) == 0 {
		return Block{}, ErrBlockListRequiresItems
	}

	if err := validateChildren(children, []BlockType{BlockTypeListItem}); err != nil {
		return Block{}, err
	}

	return Block{blockType: BlockTypeList, children: children}, nil
}

func NewListItemBlock(value string, children []Block) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	if err := validateChildren(children, []BlockType{BlockTypeItalic, BlockTypeBold, BlockTypeLink}); err != nil {
		return Block{}, err
	}

	return Block{blockType: BlockTypeListItem, value: value, children: children}, nil
}

func NewBoldBlock(value string) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	return Block{blockType: BlockTypeBold, value: value}, nil
}

func NewItalicBlock(value string) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	return Block{blockType: BlockTypeItalic, value: value}, nil
}

func NewLinkBlock(value string, children []Block) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	if len(children) != linkRequiredChildrenCount {
		return Block{}, ErrBlockLinkRequiresChild
	}

	if children[0].blockType != BlockTypeLinkURL {
		return Block{}, ErrBlockLinkRequiresChild
	}

	return Block{blockType: BlockTypeLink, value: value, children: children}, nil
}

func NewLinkURLBlock(value string) (Block, error) {
	if value == "" {
		return Block{}, ErrBlockEmptyValue
	}

	return Block{blockType: BlockTypeLinkURL, value: value}, nil
}

func (v Block) BlockType() BlockType { return v.blockType }
func (v Block) Value() string        { return v.value }
func (v Block) Children() []Block    { return v.children }

func (v Block) Equals(other Block) bool {
	if v.blockType != other.blockType || v.value != other.value {
		return false
	}

	if len(v.children) != len(other.children) {
		return false
	}

	for i := range v.children {
		if !v.children[i].Equals(other.children[i]) {
			return false
		}
	}

	return true
}

func (v Block) ToDTO() BlockDTO {
	var childrenDTO []BlockDTO
	if len(v.children) > 0 {
		childrenDTO = make([]BlockDTO, len(v.children))
		for i, child := range v.children {
			childrenDTO[i] = child.ToDTO()
		}
	}

	return BlockDTO{
		BlockType: string(v.blockType),
		Value:     v.value,
		Children:  childrenDTO,
	}
}

func validateChildren(children []Block, allowedTypes []BlockType) error {
	allowed := make(map[BlockType]bool, len(allowedTypes))
	for _, t := range allowedTypes {
		allowed[t] = true
	}

	for _, child := range children {
		if !allowed[child.blockType] {
			return ErrBlockInvalidChildren
		}
	}

	return nil
}
