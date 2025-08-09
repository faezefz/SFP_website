-- name: CreateLog :one
INSERT INTO logs (
  user_id,
  action,
  details
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CreateProjectLog :one
INSERT INTO logs (user_id, action, details, project_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListLogsByProject :many
SELECT * FROM logs
WHERE project_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


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
