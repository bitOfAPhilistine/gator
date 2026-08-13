-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetFeedById :one
SELECT *
FROM feeds
WHERE id = $1;

-- name: GetFeeds :many
SELECT *
FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: ResetFeeds :exec
DELETE FROM feeds;