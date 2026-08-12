-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted.*, feeds.name AS feed_name, users.username AS user_name
FROM inserted
INNER JOIN feeds ON inserted.feed_id = feeds.id
INNER JOIN users ON inserted.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT *
FROM feed_follows
WHERE user_id = $1;

-- name: ResetFeedFollows :exec
DELETE FROM feed_follows;