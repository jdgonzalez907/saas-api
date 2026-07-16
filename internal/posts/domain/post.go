package domain

import (
	"errors"
	"time"
)

const UnassignedPostID int64 = 0

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
)

var (
	ErrInvalidPostID                    = errors.New("invalid post identification")
	ErrInvalidPostStatus                = errors.New("invalid post status")
	ErrInvalidAuthorID                  = errors.New("invalid author identification")
	ErrDraftCannotHavePublicationDate   = errors.New("a draft post cannot have a publication date")
	ErrPublishedMustHavePublicationDate = errors.New("a published post must have a publication date")
	ErrPostNotFound                     = errors.New("the requested post was not found")
	ErrCreatingPost                     = errors.New("error creating post")
	ErrFindingPost                      = errors.New("error finding post")
	ErrChangingPost                     = errors.New("error updating post")
	ErrDeletingPost                     = errors.New("error deleting post")
	ErrFindingPosts                     = errors.New("error finding posts")
	ErrPostIDAlreadyExists              = errors.New("post ID already exists")
	ErrPostOwnershipMismatch            = errors.New("the authenticated user cannot modify another user's post")
)

func NewPostStatus(s string) (PostStatus, error) {
	status := PostStatus(s)
	switch status {
	case StatusDraft, StatusPublished:
		return status, nil
	default:
		return "", ErrInvalidPostStatus
	}
}

type Post struct {
	id                 int64
	contentInformation ContentInformation
	status             PostStatus
	createdAt          time.Time
	updatedAt          time.Time
	authorID           int64
	publishedAt        *time.Time
}

type PostDTO struct {
	ID int64 `json:"id"`
	ContentInformationDTO
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	AuthorID    int64      `json:"author_id"`
	PublishedAt *time.Time `json:"published_at"`
}

func (p *Post) ensureInvariants() error {
	if p.id < 0 {
		return ErrInvalidPostID
	}
	if p.authorID <= 0 {
		return ErrInvalidAuthorID
	}
	if p.status == StatusDraft && p.publishedAt != nil {
		return ErrDraftCannotHavePublicationDate
	}
	if p.status == StatusPublished && p.publishedAt == nil {
		return ErrPublishedMustHavePublicationDate
	}
	return nil
}

func NewPost(
	id int64,
	contentInformation ContentInformation,
	status PostStatus,
	createdAt time.Time,
	updatedAt time.Time,
	authorID int64,
	publishedAt *time.Time,
) (*Post, error) {
	if id <= 0 {
		return nil, ErrInvalidPostID
	}

	post := &Post{
		id:                 id,
		contentInformation: contentInformation,
		status:             status,
		createdAt:          createdAt.UTC(),
		updatedAt:          updatedAt.UTC(),
		authorID:           authorID,
		publishedAt:        publishedAt,
	}

	if err := post.ensureInvariants(); err != nil {
		return nil, err
	}

	return post, nil
}

func NewPostWithoutID(contentInformation ContentInformation, status PostStatus, authorID int64) (*Post, error) {
	var publishedAt *time.Time
	if status == StatusPublished {
		now := time.Now().UTC()
		publishedAt = &now
	}

	now := time.Now().UTC()
	post := &Post{
		id:                 UnassignedPostID,
		contentInformation: contentInformation,
		status:             status,
		createdAt:          now,
		updatedAt:          now,
		authorID:           authorID,
		publishedAt:        publishedAt,
	}

	if err := post.ensureInvariants(); err != nil {
		return nil, err
	}

	return post, nil
}

func (p *Post) ID() int64 {
	return p.id
}

func (p *Post) AssignID(id int64) {
	p.id = id
}

func (p *Post) ContentInformation() ContentInformation {
	return p.contentInformation
}

func (p *Post) Status() PostStatus {
	return p.status
}

func (p *Post) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Post) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Post) AuthorID() int64 {
	return p.authorID
}

func (p *Post) PublishedAt() *time.Time {
	return p.publishedAt
}

func (p *Post) IsSameAuthor(authorID int64) error {
	if p.authorID != authorID {
		return ErrPostOwnershipMismatch
	}
	return nil
}

func (p *Post) UpdateContentAndStatus(contentInformation ContentInformation, status PostStatus, authorID int64) (*Post, error) {
	if err := p.IsSameAuthor(authorID); err != nil {
		return nil, err
	}

	var publishedAt *time.Time
	if status == StatusPublished {
		if p.status == StatusPublished {
			publishedAt = p.publishedAt
		} else {
			now := time.Now().UTC()
			publishedAt = &now
		}
	}

	post := &Post{
		id:                 p.id,
		contentInformation: contentInformation,
		status:             status,
		createdAt:          p.createdAt,
		updatedAt:          time.Now().UTC(),
		authorID:           p.authorID,
		publishedAt:        publishedAt,
	}

	if err := post.ensureInvariants(); err != nil {
		return nil, err
	}

	return post, nil
}

func (p *Post) Equals(other *Post) bool {
	if other == nil {
		return false
	}
	if p.id != other.id {
		return false
	}
	if p.status != other.status {
		return false
	}
	if p.authorID != other.authorID {
		return false
	}
	if !p.createdAt.Equal(other.createdAt) {
		return false
	}
	if !p.updatedAt.Equal(other.updatedAt) {
		return false
	}
	if (p.publishedAt == nil) != (other.publishedAt == nil) {
		return false
	}
	if p.publishedAt != nil && !p.publishedAt.Equal(*other.publishedAt) {
		return false
	}
	return p.contentInformation.Equals(other.contentInformation)
}

func (p *Post) ToDTO() *PostDTO {
	var publishedAt *time.Time
	if p.publishedAt != nil {
		tVal := *p.publishedAt
		publishedAt = &tVal
	}

	return &PostDTO{
		ID:                    p.id,
		ContentInformationDTO: p.contentInformation.ToDTO(),
		Status:                string(p.status),
		CreatedAt:             p.createdAt,
		UpdatedAt:             p.updatedAt,
		AuthorID:              p.authorID,
		PublishedAt:           publishedAt,
	}
}

func PostFromDTO(dto *PostDTO) (*Post, error) {
	if dto == nil {
		return nil, nil
	}
	contentInfo, err := ContentInformationFromDTO(dto.ContentInformationDTO)
	if err != nil {
		return nil, err
	}
	status, err := NewPostStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	if dto.ID == UnassignedPostID {
		post := &Post{
			id:                 UnassignedPostID,
			contentInformation: contentInfo,
			status:             status,
			createdAt:          dto.CreatedAt.UTC(),
			updatedAt:          dto.UpdatedAt.UTC(),
			authorID:           dto.AuthorID,
			publishedAt:        dto.PublishedAt,
		}
		if err := post.ensureInvariants(); err != nil {
			return nil, err
		}
		return post, nil
	}

	return NewPost(
		dto.ID,
		contentInfo,
		status,
		dto.CreatedAt.UTC(),
		dto.UpdatedAt.UTC(),
		dto.AuthorID,
		dto.PublishedAt,
	)
}
