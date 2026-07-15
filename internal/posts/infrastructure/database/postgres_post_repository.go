package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"jdgonzalez907/saas-api/internal/posts/domain"
	"jdgonzalez907/saas-api/internal/shared/infrastructure/postgres"
)

type postgresPostRepository struct {
	queries *postgres.Queries
	pool    *pgxpool.Pool
}

func NewPostgresPostRepository(pool *pgxpool.Pool) domain.PostRepository {
	return &postgresPostRepository{
		queries: postgres.New(pool),
		pool:    pool,
	}
}

func (r *postgresPostRepository) FindByID(ctx context.Context, id int64) (*domain.Post, error) {
	row, err := r.queries.FindPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return mapRowToDomain(
		row.ID,
		row.Title,
		row.Content,
		row.Status,
		row.AuthorID,
		row.LastEditorID,
		row.PublishedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func (r *postgresPostRepository) FindAll(ctx context.Context, status domain.PostStatus, pagination domain.Pagination) ([]*domain.Post, error) {
	limit := pagination.Limit()
	var posts []*domain.Post

	if pagination.LastID() != nil {
		var lastPublishedAt pgtype.Timestamptz
		if pagination.LastPublishedAt() != nil {
			lastPublishedAt = pgtype.Timestamptz{Time: *pagination.LastPublishedAt(), Valid: true}
		}

		dbRows, err := r.queries.FindPostsPaginatedWithCursor(ctx, postgres.FindPostsPaginatedWithCursorParams{
			Status:          string(status),
			LastPublishedAt: lastPublishedAt,
			ID:              *pagination.LastID(),
			Limit:           limit,
		})
		if err != nil {
			return nil, err
		}

		posts = make([]*domain.Post, len(dbRows))
		for i, row := range dbRows {
			var err error
			posts[i], err = mapRowToDomain(
				row.ID,
				row.Title,
				row.Content,
				row.Status,
				row.AuthorID,
				row.LastEditorID,
				row.PublishedAt,
				row.CreatedAt,
				row.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
	} else {
		dbRows, err := r.queries.FindPostsPaginatedWithoutCursor(ctx, postgres.FindPostsPaginatedWithoutCursorParams{
			Status: string(status),
			Limit:  limit,
		})
		if err != nil {
			return nil, err
		}

		posts = make([]*domain.Post, len(dbRows))
		for i, row := range dbRows {
			var err error
			posts[i], err = mapRowToDomain(
				row.ID,
				row.Title,
				row.Content,
				row.Status,
				row.AuthorID,
				row.LastEditorID,
				row.PublishedAt,
				row.CreatedAt,
				row.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return posts, nil
}

func (r *postgresPostRepository) Create(ctx context.Context, post *domain.Post) error {
	contentBytes, err := toJSONB(post.ContentInformation().Content())
	if err != nil {
		return err
	}

	var publishedAt pgtype.Timestamptz
	if post.PublishedAt() != nil {
		publishedAt = pgtype.Timestamptz{Time: *post.PublishedAt(), Valid: true}
	}

	id, err := r.queries.CreatePost(ctx, postgres.CreatePostParams{
		Title:        post.ContentInformation().Title(),
		Content:      contentBytes,
		Status:       string(post.Status()),
		AuthorID:     post.AuthorID(),
		LastEditorID: post.LastEditorID(),
		PublishedAt:  publishedAt,
		CreatedAt:    pgtype.Timestamptz{Time: post.CreatedAt(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: post.UpdatedAt(), Valid: true},
	})
	if err != nil {
		return err
	}

	post.AssignID(id)
	return nil
}

func (r *postgresPostRepository) Update(ctx context.Context, post *domain.Post) error {
	contentBytes, err := toJSONB(post.ContentInformation().Content())
	if err != nil {
		return err
	}

	var publishedAt pgtype.Timestamptz
	if post.PublishedAt() != nil {
		publishedAt = pgtype.Timestamptz{Time: *post.PublishedAt(), Valid: true}
	}

	return r.queries.UpdatePost(ctx, postgres.UpdatePostParams{
		Title:        post.ContentInformation().Title(),
		Content:      contentBytes,
		Status:       string(post.Status()),
		AuthorID:     post.AuthorID(),
		LastEditorID: post.LastEditorID(),
		PublishedAt:  publishedAt,
		UpdatedAt:    pgtype.Timestamptz{Time: post.UpdatedAt(), Valid: true},
		ID:           post.ID(),
	})
}

func (r *postgresPostRepository) Delete(ctx context.Context, id int64, deletedByID int64) error {
	return r.queries.DeletePost(ctx, postgres.DeletePostParams{
		ID:        id,
		DeletedBy: pgtype.Int8{Int64: deletedByID, Valid: true},
	})
}

func toJSONB(blocks []domain.Block) ([]byte, error) {
	if blocks == nil {
		return nil, nil
	}

	blocksDTO := make([]domain.BlockDTO, len(blocks))
	for i, block := range blocks {
		blocksDTO[i] = block.ToDTO()
	}

	bytes, err := json.Marshal(blocksDTO)
	if err != nil {
		return nil, fmt.Errorf("error marshaling blocks: %w", err)
	}
	return bytes, nil
}

func mapRowToDomain(
	id int64,
	title string,
	content []byte,
	status string,
	authorID int64,
	lastEditorID int64,
	publishedAt pgtype.Timestamptz,
	createdAt pgtype.Timestamptz,
	updatedAt pgtype.Timestamptz,
) (*domain.Post, error) {
	var blocksDTO []domain.BlockDTO
	if err := json.Unmarshal(content, &blocksDTO); err != nil {
		return nil, fmt.Errorf("error deserializing post content: %w", err)
	}

	contentInfo, err := domain.ContentInformationFromDTO(domain.ContentInformationDTO{
		Title:   title,
		Content: blocksDTO,
	})
	if err != nil {
		return nil, err
	}

	statusVO, err := domain.NewPostStatus(status)
	if err != nil {
		return nil, err
	}

	var pubAt *time.Time
	if publishedAt.Valid {
		t := publishedAt.Time.UTC()
		pubAt = &t
	}

	return domain.NewPost(domain.PostParams{
		ID:                 id,
		ContentInformation: contentInfo,
		Status:             statusVO,
		CreatedAt:          createdAt.Time.UTC(),
		UpdatedAt:          updatedAt.Time.UTC(),
		AuthorID:           authorID,
		LastEditorID:       lastEditorID,
		PublishedAt:        pubAt,
	})
}
