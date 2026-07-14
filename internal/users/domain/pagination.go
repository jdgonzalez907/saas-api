package domain

import "errors"

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
	lastID *int64
	limit  int32
}

func NewPagination(lastID *int64, limitVal *int32) (Pagination, error) {
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

func (p Pagination) LastID() *int64 {
	return p.lastID
}

func (p Pagination) Limit() int32 {
	return p.limit
}
