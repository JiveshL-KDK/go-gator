
-- name: CreateNewPost :one
INSERT INTO posts (id, created_at, updated_at, published_at, title, url, description, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: CreatePosts :many
INSERT INTO posts (id, created_at, updated_at, published_at, title, url, description, feed_id)
SELECT
    unnest(sqlc.arg(ids)::uuid[]),
    unnest(sqlc.arg(created_ats)::timestamp[]),
    unnest(sqlc.arg(updated_ats)::timestamp[]),
    unnest(sqlc.arg(published_ats)::timestamp[]),
    unnest(sqlc.arg(titles)::text[]),
    unnest(sqlc.arg(urls)::text[]),
    unnest(sqlc.arg(descriptions)::text[]),
    unnest(sqlc.arg(feed_ids)::uuid[])
RETURNING *;

-- name: GetPostsForAUser :many
WITH feeds_user_follows AS (
    SELECT * FROM feed_follows
    WHERE user_id = $1
)
SELECT posts.* FROM posts
INNER JOIN feeds_user_follows ON posts.feed_id = feeds_user_follows.feed_id
LIMIT $2;