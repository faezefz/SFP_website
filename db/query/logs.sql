-- name: CreateLog :one
INSERT INTO logs (user_id, project_id, action, details)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLogByID :one
SELECT * FROM logs WHERE id = $1 LIMIT 1;

-- name: GetLogsByProjectOrUser :many
SELECT * FROM logs
WHERE project_id = $1 OR user_id = $2
ORDER BY created_at DESC;

-- name: DeleteLog :exec
DELETE FROM logs WHERE id = $1;
