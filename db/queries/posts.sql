-- name: CreatePost :one
INSERT INTO posts (
    title,
    content,
    status,
    author_id,
    published_at,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdatePost :exec
UPDATE posts SET
    title = $1,
    content = $2,
    status = $3,
    published_at = $4,
    updated_at = $5
WHERE id = $6 AND deleted_at IS NULL;

-- name: DeletePost :exec
UPDATE posts SET
    deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPostByID :one
SELECT id, title, content, status, author_id, published_at, created_at, updated_at
FROM posts
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindPostsPaginatedWithoutCursor :many
SELECT id, title, content, status, author_id, published_at, created_at, updated_at
FROM posts
WHERE status = $1 AND author_id = $3 AND deleted_at IS NULL
ORDER BY published_at DESC, id DESC
LIMIT $2;

-- name: FindPostsPaginatedWithCursor :many
SELECT id, title, content, status, author_id, published_at, created_at, updated_at
FROM posts
WHERE status = $1 AND author_id = $4 AND deleted_at IS NULL AND (
    (sqlc.narg('last_published_at')::timestamptz IS NOT NULL AND published_at < sqlc.narg('last_published_at')) OR 
    (sqlc.narg('last_published_at')::timestamptz IS NOT NULL AND published_at = sqlc.narg('last_published_at') AND id < $2) OR
    (sqlc.narg('last_published_at')::timestamptz IS NULL AND id < $2)
)
ORDER BY published_at DESC, id DESC
LIMIT $3;
