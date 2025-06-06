-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feeds_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowByIDS :one
SELECT * FROM feeds_follows WHERE user_id = $1 AND feed_id = $2;

-- name: GetFeedFollowsForUser :many
SELECT 
    feeds_follows.*, 
    feeds.name AS feed_name 
FROM feeds_follows 
INNER JOIN feeds ON feeds.id = feeds_follows.feed_id
WHERE feeds_follows.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feeds_follows
WHERE user_id = $1 AND feed_id = $2;