package domain

type Pagination struct {
	lastID *int
	limit  int
}

func NewPagination(lastID *int, limitVal int) Pagination {
	limit := 10
	if limitVal == 10 || limitVal == 25 || limitVal == 50 {
		limit = limitVal
	}
	return Pagination{
		lastID: lastID,
		limit:  limit,
	}
}

func (p Pagination) LastID() *int {
	return p.lastID
}

func (p Pagination) Limit() int {
	return p.limit
}
