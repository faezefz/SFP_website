-- name: CreateModel :one
INSERT INTO models (
  user_id,
  name,
  description,
  model_type,
  file_path
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetModel :one
SELECT * FROM models
WHERE id = $1 LIMIT 1;

-- name: ListModels :many
SELECT * FROM models
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateModel :one
UPDATE models
  SET name = $2,
      description = $3,
      model_type = $4,
      file_path = $5
WHERE id = $1
RETURNING *;

-- name: DeleteModel :exec
DELETE FROM models
WHERE id = $1;
