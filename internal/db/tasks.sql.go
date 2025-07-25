// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: tasks.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (user_id, description)
VALUES ($1, $2)
RETURNING id, user_id, description, created_at
`

type CreateTaskParams struct {
	UserID      pgtype.UUID `json:"user_id"`
	Description string      `json:"description"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRow(ctx, createTask, arg.UserID, arg.Description)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const readTask = `-- name: ReadTask :one
SELECT id, user_id, description, created_at
FROM tasks
WHERE id = $1
`

func (q *Queries) ReadTask(ctx context.Context, id pgtype.UUID) (Task, error) {
	row := q.db.QueryRow(ctx, readTask, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const removeTask = `-- name: RemoveTask :exec
DELETE FROM tasks
WHERE id = $1
`

func (q *Queries) RemoveTask(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, removeTask, id)
	return err
}

const updateTask = `-- name: UpdateTask :one
UPDATE tasks
SET description = $2
WHERE id = $1
RETURNING id, user_id, description, created_at
`

type UpdateTaskParams struct {
	ID          pgtype.UUID `json:"id"`
	Description string      `json:"description"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error) {
	row := q.db.QueryRow(ctx, updateTask, arg.ID, arg.Description)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}
