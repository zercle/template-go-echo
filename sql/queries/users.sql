-- name: CreateUser :exec
INSERT INTO users (id, name, email)
VALUES (?, ?, ?);

-- name: GetUserByID :one
SELECT id, name, email, created_at, updated_at
FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT id, name, email, created_at, updated_at
FROM users
WHERE email = ?;

-- name: UpdateUser :exec
UPDATE users
SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT id, name, email, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsers :one
SELECT COUNT(*) as total
FROM users;
