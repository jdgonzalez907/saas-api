-- name: CreateUser :one
INSERT INTO users (
    identification_type,
    identification_number,
    first_name,
    last_name,
    birth_date,
    address,
    phone_country_code,
    phone_number,
    email,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users SET 
    identification_type = $1, 
    identification_number = $2, 
    first_name = $3, 
    last_name = $4, 
    birth_date = $5, 
    address = $6, 
    phone_country_code = $7, 
    phone_number = $8, 
    email = $9, 
    updated_at = $10 
WHERE id = $11 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users SET 
    deleted_at = NOW(), 
    updated_at = NOW() 
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindUserByID :one
SELECT 
    id, identification_type, identification_number, first_name, last_name, 
    birth_date, address, phone_country_code, phone_number, email, 
    created_at, updated_at 
FROM users 
WHERE id = $1 AND deleted_at IS NULL;

-- name: FindUserByPhone :one
SELECT 
    id, identification_type, identification_number, first_name, last_name, 
    birth_date, address, phone_country_code, phone_number, email, 
    created_at, updated_at 
FROM users 
WHERE phone_country_code = $1 AND phone_number = $2 AND deleted_at IS NULL;

-- name: FindUserByEmail :one
SELECT 
    id, identification_type, identification_number, first_name, last_name, 
    birth_date, address, phone_country_code, phone_number, email, 
    created_at, updated_at 
FROM users 
WHERE email = $1 AND deleted_at IS NULL;

-- name: FindUsersPaginatedWithCursor :many
SELECT 
    id, identification_type, identification_number, first_name, last_name, 
    birth_date, address, phone_country_code, phone_number, email, 
    created_at, updated_at 
FROM users 
WHERE id > $1 AND deleted_at IS NULL 
ORDER BY id ASC 
LIMIT $2;

-- name: FindUsersPaginatedWithoutCursor :many
SELECT 
    id, identification_type, identification_number, first_name, last_name, 
    birth_date, address, phone_country_code, phone_number, email, 
    created_at, updated_at 
FROM users 
WHERE deleted_at IS NULL 
ORDER BY id ASC 
LIMIT $1;
