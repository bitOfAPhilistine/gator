-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT *
FROM users;

-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;