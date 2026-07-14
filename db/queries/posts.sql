-- name: CreatePost :one
INSERT INTO posts (
    title,
    content,
    status,
    author_id,
    last_editor_id,
    published_at,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: UpdatePost :exec
UPDATE posts SET
    title = $1,
    content = $2,
    status = $3,
    author_id = $4,
    last_editor_id = $5,
    published_at = $6,
    updated_at = $7
WHERE id = $8 AND deleted_at IS NULL;

-- name: DeletePost :exec
UPDATE posts SET
    deleted_at = NOW(),
    deleted_by = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPostByID :one
SELECT id, title, content, status, author_id, last_editor_id, published_at, created_at, updated_at
FROM posts
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPostsPaginatedWithoutCursor :many
SELECT id, title, content, status, author_id, last_editor_id, published_at, created_at, updated_at
FROM posts
WHERE status = $1 AND deleted_at IS NULL
ORDER BY published_at DESC, id DESC
LIMIT $2;

-- name: FindPostsPaginatedWithCursor :many
SELECT id, title, content, status, author_id, last_editor_id, published_at, created_at, updated_at
FROM posts
WHERE status = $1 AND deleted_at IS NULL AND (
    (sqlc.narg('last_published_at')::timestamptz IS NOT NULL AND published_at < sqlc.narg('last_published_at')) OR 
    (sqlc.narg('last_published_at')::timestamptz IS NOT NULL AND published_at = sqlc.narg('last_published_at') AND id < $2) OR
    (sqlc.narg('last_published_at')::timestamptz IS NULL AND id < $2)
)
ORDER BY published_at DESC, id DESC
LIMIT $3;
