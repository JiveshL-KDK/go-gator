-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsOfUser :many
SELECT * FROM feeds
WHERE user_id = $1;


-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: UserHasFeed :one
SELECT EXISTS(
    SELECT 1 FROM feeds
    WHERE user_id = $1
    AND url = $2
) as exists;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: RemoveAllFeeds :exec
TRUNCATE TABLE feeds;