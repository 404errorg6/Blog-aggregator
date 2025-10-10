-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
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

SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds 
ON inserted_feed_follow.user_id = feeds.id
INNER JOIN users
ON inserted_feed_follow.feed_id = users.id
ORDER BY user_name;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id,
        f.name AS feed_name,
        u.name AS user_name
  FROM feed_follows AS ff
INNER JOIN feeds AS f
ON ff.feed_id = f.id
INNER JOIN users AS u
ON ff.user_id = u.id
WHERE ff.user_id = $1;

-- name: DelFeedFollowEntry :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
