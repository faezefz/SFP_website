-- name: CreateProject :one
INSERT INTO projects (owner_user_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = $1 LIMIT 1;

-- name: GetProjectsByOwnerID :many
SELECT * FROM projects WHERE owner_user_id = $1 ORDER BY id;

-- name: UpdateProject :one
UPDATE projects
SET name = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;
