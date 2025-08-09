-- name: CreateProject :one
INSERT INTO projects (owner_user_id, name, description, visibility)
VALUES ($1, $2, $3, COALESCE($4, 'private'))
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1 LIMIT 1;

-- name: ListProjectsByOwner :many
SELECT * FROM projects
WHERE owner_user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProject :one
UPDATE projects
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    visibility = COALESCE($4, visibility)
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;
