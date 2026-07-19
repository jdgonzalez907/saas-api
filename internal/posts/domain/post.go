package domain

import (
	"errors"
	"time"
)

var (
	ErrPostTitleRequired      = errors.New("post title is required")
	ErrPostSlugRequired       = errors.New("post slug is required")
	ErrPostCoverRequired      = errors.New("post cover is required")
	ErrPostStatusRequired     = errors.New("post status is required")
	ErrPostAuthorIDRequired   = errors.New("post author ID is required")
	ErrPostIDRequired         = errors.New("post ID is required")
	ErrPostInvalidRootBlock   = errors.New("post content contains invalid root block type")
	ErrPostUnauthorizedUpdate = errors.New("post cannot be updated by this user")
	ErrPostSlugAlreadyExists  = errors.New("post slug already exists")
	ErrPostNotFound           = errors.New("post not found")

	ErrCreatePost   = errors.New("cannot create post")
	ErrUpdateContent = errors.New("cannot update content")
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

type Post struct {
	id        int64
	title     string
	slug      string
	cover     string
	content   []Block
	status    PostStatus
	authorID  int64
	createdAt time.Time
	updatedAt time.Time
}

type PostDTO struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Cover     string     `json:"cover"`
	Content   []BlockDTO `json:"content"`
	Status    string     `json:"status"`
	AuthorID  int64      `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func New(title, slug, cover string, content []Block, status PostStatus, authorID int64) (*Post, error) {
	if err := validatePostTitle(title); err != nil {
		return nil, err
	}
	if err := validatePostSlug(slug); err != nil {
		return nil, err
	}
	if err := validatePostCover(cover); err != nil {
		return nil, err
	}
	if err := validatePostStatus(status); err != nil {
		return nil, err
	}
	if err := validatePostAuthorID(authorID); err != nil {
		return nil, err
	}
	if err := validateRootBlocks(content); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Post{
		title:     title,
		slug:      slug,
		cover:     cover,
		content:   content,
		status:    status,
		authorID:  authorID,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func NewWithID(id int64, title, slug, cover string, content []Block, status PostStatus, authorID int64, createdAt, updatedAt time.Time) (*Post, error) {
	if id <= 0 {
		return nil, ErrPostIDRequired
	}
	if err := validatePostTitle(title); err != nil {
		return nil, err
	}
	if err := validatePostSlug(slug); err != nil {
		return nil, err
	}
	if err := validatePostCover(cover); err != nil {
		return nil, err
	}
	if err := validatePostStatus(status); err != nil {
		return nil, err
	}
	if err := validatePostAuthorID(authorID); err != nil {
		return nil, err
	}
	if err := validateRootBlocks(content); err != nil {
		return nil, err
	}

	return &Post{
		id:        id,
		title:     title,
		slug:      slug,
		cover:     cover,
		content:   content,
		status:    status,
		authorID:  authorID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (p *Post) ID() int64            { return p.id }
func (p *Post) Title() string        { return p.title }
func (p *Post) Slug() string         { return p.slug }
func (p *Post) Cover() string        { return p.cover }
func (p *Post) Content() []Block     { return p.content }
func (p *Post) Status() PostStatus   { return p.status }
func (p *Post) AuthorID() int64      { return p.authorID }
func (p *Post) CreatedAt() time.Time { return p.createdAt }
func (p *Post) UpdatedAt() time.Time { return p.updatedAt }

func (p *Post) AssignID(id int64) {
	p.id = id
}

func (p *Post) Equals(other *Post) bool {
	if other == nil {
		return false
	}
	return p.id == other.id
}

func (p *Post) UpdateContent(title, slug, cover string, content []Block, status PostStatus, executedBy int64) error {
	if executedBy != p.authorID {
		return ErrPostUnauthorizedUpdate
	}
	if err := validatePostTitle(title); err != nil {
		return err
	}
	if err := validatePostSlug(slug); err != nil {
		return err
	}
	if err := validatePostCover(cover); err != nil {
		return err
	}
	if err := validatePostStatus(status); err != nil {
		return err
	}
	if err := validateRootBlocks(content); err != nil {
		return err
	}

	p.title = title
	p.slug = slug
	p.cover = cover
	p.content = content
	p.status = status
	p.updatedAt = time.Now()
	return nil
}

func (p *Post) ToDTO() PostDTO {
	var contentDTO []BlockDTO
	if len(p.content) > 0 {
		contentDTO = make([]BlockDTO, len(p.content))
		for i, block := range p.content {
			contentDTO[i] = block.ToDTO()
		}
	}

	return PostDTO{
		ID:        p.id,
		Title:     p.title,
		Slug:      p.slug,
		Cover:     p.cover,
		Content:   contentDTO,
		Status:    string(p.status),
		AuthorID:  p.authorID,
		CreatedAt: p.createdAt,
		UpdatedAt: p.updatedAt,
	}
}

func validatePostTitle(title string) error {
	if title == "" {
		return ErrPostTitleRequired
	}
	return nil
}

func validatePostSlug(slug string) error {
	if slug == "" {
		return ErrPostSlugRequired
	}
	return nil
}

func validatePostCover(cover string) error {
	if cover == "" {
		return ErrPostCoverRequired
	}
	return nil
}

func validatePostStatus(status PostStatus) error {
	if status == "" {
		return ErrPostStatusRequired
	}
	return nil
}

func validatePostAuthorID(authorID int64) error {
	if authorID <= 0 {
		return ErrPostAuthorIDRequired
	}
	return nil
}

func validateRootBlocks(content []Block) error {
	allowed := map[BlockType]bool{
		BlockTypeTitle:     true,
		BlockTypeSubtitle:  true,
		BlockTypeParagraph: true,
		BlockTypeList:      true,
		BlockTypeImage:     true,
	}

	for _, block := range content {
		if !allowed[block.BlockType()] {
			return ErrPostInvalidRootBlock
		}
	}
	return nil
}
