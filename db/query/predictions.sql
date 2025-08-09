-- name: CreatePrediction :one
INSERT INTO predictions (
  user_id,
  dataset_id,
  model_id,
  result_file_path,
  status,
  project_id         
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetPrediction :one
SELECT * FROM predictions
WHERE id = $1 LIMIT 1;

-- name: ListPredictions :many
SELECT * FROM predictions
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdatePrediction :one
UPDATE predictions
  SET result_file_path = $2,
      status = $3
WHERE id = $1
RETURNING *;

-- name: DeletePrediction :exec
DELETE FROM predictions
WHERE id = $1;

-- name: ListPredictionsByProject :many
SELECT * FROM predictions
WHERE project_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPredictionsByProjectAndStatus :many
SELECT * FROM predictions
WHERE project_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
