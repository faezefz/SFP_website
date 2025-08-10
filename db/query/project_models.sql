-- name: AddModelToProject :exec
INSERT INTO project_models (project_id, model_id)
VALUES ($1, $2);

-- name: GetModelsByProjectID :many
SELECT m.*
FROM models m
JOIN project_models pm ON m.id = pm.model_id
WHERE pm.project_id = $1;

-- name: RemoveModelFromProject :exec
DELETE FROM project_models
WHERE project_id = $1 AND model_id = $2;
