-- name: CreateDataset :one
INSERT INTO datasets (
  user_id,
  name,
  description,
  file_path
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetDataset :one
SELECT * FROM datasets
WHERE id = $1 LIMIT 1;

-- name: ListDatasets :many
SELECT * FROM datasets
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListDatasetsByUserID :many
SELECT * FROM datasets
WHERE user_id = $1
ORDER BY uploaded_at DESC;

-- name: UpdateDataset :one
UPDATE datasets
  SET name = $2,
      description = $3,
      file_path = $4
WHERE id = $1
RETURNING *;

-- name: DeleteDataset :exec
DELETE FROM datasets
WHERE id = $1;
