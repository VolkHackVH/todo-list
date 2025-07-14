-- name: CreateUser :one
INSERT INTO users (username)
VALUES ($1)
RETURNING *;

-- name: ReadUserInfo :one
SELECT *
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING *;

-- name: RemoveUser :exec
DELETE FROM users
WHERE id = $1;
