package domain

import (
	"errors"
)

var (
	ErrEmptyPostTitle = errors.New("post title cannot be empty")
	ErrOrphanBlock    = errors.New("block cannot exist without a valid parent")
)

type ContentInformation struct {
	title   string
	content []Block
}

type ContentInformationDTO struct {
	Title   string     `json:"title"`
	Content []BlockDTO `json:"content"`
}

func NewContentInformation(title string, content []Block) (ContentInformation, error) {
	if title == "" {
		return ContentInformation{}, ErrEmptyPostTitle
	}

	for _, block := range content {
		if block.blockType == TypeMark || block.blockType == TypeListItem {
			return ContentInformation{}, ErrOrphanBlock
		}
	}

	return ContentInformation{
		title:   title,
		content: content,
	}, nil
}

func (c ContentInformation) Title() string {
	return c.title
}

func (c ContentInformation) Content() []Block {
	return c.content
}

func (c ContentInformation) Equals(other ContentInformation) bool {
	if c.title != other.title {
		return false
	}
	if len(c.content) != len(other.content) {
		return false
	}
	for i := range c.content {
		if !c.content[i].Equals(other.content[i]) {
			return false
		}
	}
	return true
}

func (c ContentInformation) ToDTO() ContentInformationDTO {
	var contentDTO []BlockDTO
	if c.content != nil {
		contentDTO = make([]BlockDTO, len(c.content))
		for i, block := range c.content {
			contentDTO[i] = block.ToDTO()
		}
	}
	return ContentInformationDTO{
		Title:   c.title,
		Content: contentDTO,
	}
}

func ContentInformationFromDTO(dto ContentInformationDTO) (ContentInformation, error) {
	var content []Block
	if dto.Content != nil {
		content = make([]Block, len(dto.Content))
		for i, blockDTO := range dto.Content {
			block, err := BlockFromDTO(blockDTO)
			if err != nil {
				return ContentInformation{}, err
			}
			content[i] = block
		}
	}
	return NewContentInformation(dto.Title, content)
}
