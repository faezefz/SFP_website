-- name: AddModelToProject :exec
INSERT INTO project_models (project_id, model_id)
VALUES ($1, $2)
ON CONFLICT (project_id, model_id) DO NOTHING;

-- name: RemoveModelFromProject :exec
DELETE FROM project_models WHERE project_id = $1 AND model_id = $2;

-- name: ListProjectModels :many
SELECT m.*
FROM project_models pm
JOIN models m ON m.id = pm.model_id
WHERE pm.project_id = $1
ORDER BY m.created_at DESC;
