-- SQL queries for user domain

-- name: CreateUser :exec
INSERT INTO users (id, email, name, password_hash, is_active, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, NOW(), NOW());

-- name: GetUserByID :one
SELECT id, email, name, password_hash, is_active, created_at, updated_at, deleted_at
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, is_active, created_at, updated_at, deleted_at
FROM users
WHERE email = ? AND deleted_at IS NULL;

-- name: UpdateUser :exec
UPDATE users
SET name = ?, email = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, email, name, password_hash, is_active, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetUserCount :one
SELECT COUNT(*) as count
FROM users
WHERE deleted_at IS NULL;
