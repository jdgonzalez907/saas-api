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
	ErrInvalidPostID     = errors.New("invalid post id")
	ErrInvalidPostStatus = errors.New("invalid post status")
	ErrInvalidUserID     = errors.New("invalid user id")
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
	createdBy          int64
	updatedBy          int64
}

type PostParams struct {
	ID                 int64
	ContentInformation ContentInformation
	Status             PostStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CreatedBy          int64
	UpdatedBy          int64
}

type PostDTO struct {
	ID int64 `json:"id"`
	ContentInformationDTO
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

func NewPost(params PostParams) (*Post, error) {
	if params.ID <= 0 {
		return nil, ErrInvalidPostID
	}
	if params.CreatedBy <= 0 || params.UpdatedBy <= 0 {
		return nil, ErrInvalidUserID
	}

	return &Post{
		id:                 params.ID,
		contentInformation: params.ContentInformation,
		status:             params.Status,
		createdAt:          params.CreatedAt.UTC(),
		updatedAt:          params.UpdatedAt.UTC(),
		createdBy:          params.CreatedBy,
		updatedBy:          params.UpdatedBy,
	}, nil
}

func NewPostWithoutID(contentInformation ContentInformation, status PostStatus, createdBy int64) (*Post, error) {
	if createdBy <= 0 {
		return nil, ErrInvalidUserID
	}

	now := time.Now().UTC()
	return &Post{
		id:                 UnassignedPostID,
		contentInformation: contentInformation,
		status:             status,
		createdAt:          now,
		updatedAt:          now,
		createdBy:          createdBy,
		updatedBy:          createdBy,
	}, nil
}

func (p *Post) ID() int64 {
	return p.id
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

func (p *Post) CreatedBy() int64 {
	return p.createdBy
}

func (p *Post) UpdatedBy() int64 {
	return p.updatedBy
}

func (p *Post) WithContentAndStatus(contentInformation ContentInformation, status PostStatus, updatedBy int64) (*Post, error) {
	if updatedBy <= 0 {
		return nil, ErrInvalidUserID
	}

	return &Post{
		id:                 p.id,
		contentInformation: contentInformation,
		status:             status,
		createdAt:          p.createdAt,
		updatedAt:          time.Now().UTC(),
		createdBy:          p.createdBy,
		updatedBy:          updatedBy,
	}, nil
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
	if p.createdBy != other.createdBy {
		return false
	}
	if p.updatedBy != other.updatedBy {
		return false
	}
	if !p.createdAt.Equal(other.createdAt) {
		return false
	}
	if !p.updatedAt.Equal(other.updatedAt) {
		return false
	}
	return p.contentInformation.Equals(other.contentInformation)
}

func (p *Post) ToDTO() *PostDTO {
	return &PostDTO{
		ID:                    p.id,
		ContentInformationDTO: p.contentInformation.ToDTO(),
		Status:                string(p.status),
		CreatedAt:             p.createdAt,
		UpdatedAt:             p.updatedAt,
		CreatedBy:             p.createdBy,
		UpdatedBy:             p.updatedBy,
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
		return &Post{
			id:                 UnassignedPostID,
			contentInformation: contentInfo,
			status:             status,
			createdAt:          dto.CreatedAt.UTC(),
			updatedAt:          dto.UpdatedAt.UTC(),
			createdBy:          dto.CreatedBy,
			updatedBy:          dto.UpdatedBy,
		}, nil
	}

	return NewPost(PostParams{
		ID:                 dto.ID,
		ContentInformation: contentInfo,
		Status:             status,
		CreatedAt:          dto.CreatedAt.UTC(),
		UpdatedAt:          dto.UpdatedAt.UTC(),
		CreatedBy:          dto.CreatedBy,
		UpdatedBy:          dto.UpdatedBy,
	})
}
