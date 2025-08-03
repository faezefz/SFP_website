-- name: CreateLog :one
INSERT INTO logs (
  user_id,
  action,
  details
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetLog :one
SELECT * FROM logs
WHERE id = $1 LIMIT 1;

-- name: ListLogs :many
SELECT * FROM logs
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteLog :exec
DELETE FROM logs
WHERE id = $1;
