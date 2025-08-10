-- name: CreateDataset :one
INSERT INTO datasets (user_id, name, description, file_path)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDatasetByID :one
SELECT * FROM datasets WHERE id = $1 LIMIT 1;

-- name: GetDatasetsByUserID :many
SELECT * FROM datasets WHERE user_id = $1 ORDER BY id;

-- name: UpdateDataset :one
UPDATE datasets
SET name = $2,
    description = $3,
    file_path = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDataset :exec
DELETE FROM datasets WHERE id = $1;
