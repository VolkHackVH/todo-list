-- name: CreateTask :one
INSERT INTO tasks (user_id, description)
VALUES ($1, $2)
RETURNING *;

-- name: ReadTask :one
SELECT *
FROM tasks
WHERE id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET description = $2
WHERE id = $1
RETURNING *;

-- name: RemoveTask :exec
DELETE FROM tasks
WHERE id = $1;
