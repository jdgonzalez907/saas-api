package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidPaginationLimit  = errors.New("invalid pagination page size limit")
	ErrInvalidPaginationCursor = errors.New("invalid pagination cursor format or parameters")
)

const (
	Limit10      int32 = 10
	Limit25      int32 = 25
	Limit50      int32 = 50
	DefaultLimit int32 = Limit10
)

var AllowedLimits = map[int32]bool{
	Limit10: true,
	Limit25: true,
	Limit50: true,
}

type Pagination struct {
	lastPublishedAt *time.Time
	lastID          *int64
	limit           int32
}

func NewPagination(lastPublishedAt *time.Time, lastID *int64, limitVal *int32) (Pagination, error) {
	limit := DefaultLimit
	if limitVal != nil {
		if !AllowedLimits[*limitVal] {
			return Pagination{}, ErrInvalidPaginationLimit
		}
		limit = *limitVal
	}

	// Composite cursor: either both parameters are provided, or both are nil.
	if (lastPublishedAt != nil && lastID == nil) || (lastPublishedAt == nil && lastID != nil) {
		return Pagination{}, ErrInvalidPaginationCursor
	}

	if lastID != nil && *lastID <= UnassignedPostID {
		return Pagination{}, ErrInvalidPaginationCursor
	}

	return Pagination{
		lastPublishedAt: lastPublishedAt,
		lastID:          lastID,
		limit:           limit,
	}, nil
}

func (p Pagination) LastPublishedAt() *time.Time {
	return p.lastPublishedAt
}

func (p Pagination) LastID() *int64 {
	return p.lastID
}

func (p Pagination) Limit() int32 {
	return p.limit
}
