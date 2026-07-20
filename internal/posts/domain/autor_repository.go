package domain

import "context"

type AutorRepository interface {
	FindByID(ctx context.Context, id int64) (*Autor, error)
}
