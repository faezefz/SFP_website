-- name: CreatePrediction :one
INSERT INTO predictions (user_id, dataset_id, model_id, project_id, result_file_path)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPredictionByID :one
SELECT * FROM predictions WHERE id = $1 LIMIT 1;

-- name: GetPredictionsByUserID :many
SELECT * FROM predictions WHERE user_id = $1 ORDER BY id;

-- name: UpdatePrediction :one
UPDATE predictions
SET result_file_path = $2
WHERE id = $1
RETURNING *;

-- name: DeletePrediction :exec
DELETE FROM predictions WHERE id = $1;
