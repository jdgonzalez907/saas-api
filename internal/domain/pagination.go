package domain

import "errors"

var (
	ErrInvalidPaginationLimit  = errors.New("invalid pagination limit")
	ErrInvalidPaginationCursor = errors.New("invalid pagination cursor")
)

const (
	Limit10      = 10
	Limit25      = 25
	Limit50      = 50
	DefaultLimit = Limit10
)

var AllowedLimits = map[int]bool{
	Limit10: true,
	Limit25: true,
	Limit50: true,
}

type Pagination struct {
	lastID *int
	limit  int
}

func NewPagination(lastID *int, limitVal *int) (Pagination, error) {
	limit := DefaultLimit
	if limitVal != nil {
		if !AllowedLimits[*limitVal] {
			return Pagination{}, ErrInvalidPaginationLimit
		}
		limit = *limitVal
	}

	if lastID != nil && *lastID <= UnassignedUserID {
		return Pagination{}, ErrInvalidPaginationCursor
	}

	return Pagination{
		lastID: lastID,
		limit:  limit,
	}, nil
}

func (p Pagination) LastID() *int {
	return p.lastID
}

func (p Pagination) Limit() int {
	return p.limit
}
