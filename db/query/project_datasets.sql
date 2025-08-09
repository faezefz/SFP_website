-- name: AddDatasetToProject :exec
INSERT INTO project_datasets (project_id, dataset_id)
VALUES ($1, $2)
ON CONFLICT (project_id, dataset_id) DO NOTHING;

-- name: RemoveDatasetFromProject :exec
DELETE FROM project_datasets WHERE project_id = $1 AND dataset_id = $2;

-- name: ListProjectDatasets :many
SELECT d.*
FROM project_datasets pd
JOIN datasets d ON d.id = pd.dataset_id
WHERE pd.project_id = $1
ORDER BY d.uploaded_at DESC;
