-- SQL queries for user session domain

-- name: CreateSession :exec
INSERT INTO user_sessions (id, user_id, refresh_token_hash, ip_address, user_agent, expires_at, created_at)
VALUES (?, ?, ?, ?, ?, ?, NOW());

-- name: GetSessionByID :one
SELECT id, user_id, refresh_token_hash, ip_address, user_agent, expires_at, created_at
FROM user_sessions
WHERE id = ? AND expires_at > NOW();

-- name: GetSessionByUserID :many
SELECT id, user_id, refresh_token_hash, ip_address, user_agent, expires_at, created_at
FROM user_sessions
WHERE user_id = ? AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: DeleteSession :exec
DELETE FROM user_sessions
WHERE id = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM user_sessions
WHERE expires_at <= NOW();

-- name: GetSessionByTokenHash :one
SELECT id, user_id, refresh_token_hash, ip_address, user_agent, expires_at, created_at
FROM user_sessions
WHERE refresh_token_hash = ? AND expires_at > NOW();
