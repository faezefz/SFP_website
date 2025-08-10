-- name: AddDatasetToProject :exec
INSERT INTO project_datasets (project_id, dataset_id)
VALUES ($1, $2);

-- name: GetDatasetsByProjectID :many
SELECT d.*
FROM datasets d
JOIN project_datasets pd ON d.id = pd.dataset_id
WHERE pd.project_id = $1;

-- name: RemoveDatasetFromProject :exec
DELETE FROM project_datasets
WHERE project_id = $1 AND dataset_id = $2;
