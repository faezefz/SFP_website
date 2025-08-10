-- name: CreateModel :one
INSERT INTO models (user_id, name, description, file_path)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetModelByID :one
SELECT * FROM models WHERE id = $1 LIMIT 1;

-- name: GetModelsByUserID :many
SELECT * FROM models WHERE user_id = $1 ORDER BY id;

-- name: UpdateModel :one
UPDATE models
SET name = $2,
    description = $3,
    file_path = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteModel :exec
DELETE FROM models WHERE id = $1;
