-- name: GetUser :one
SELECT id, first_name, last_name, email, password, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, first_name, last_name, email, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, first_name, last_name, email, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, first_name, last_name, email, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

