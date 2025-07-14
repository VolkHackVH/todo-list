package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/VolkHackVH/todo-list.git/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

type TaskService struct {
	db *db.Queries
}

func NewTaskService(query *db.Queries) *TaskService {
	return &TaskService{db: query}
}

// ? Getting user_id can be done via - JWT token
func (s *TaskService) CreateTask(ctx context.Context, user_id pgtype.UUID, text string) (db.Task, error) {
	log.WithFields(log.Fields{
		"user_id": user_id,
		"text":    text,
	}).Info("Updating task")

	task, err := s.db.CreateTask(ctx, db.CreateTaskParams{
		UserID:      user_id,
		Description: text,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Task{}, fmt.Errorf("[SERVICE] user not found: %w", err)
		}
		return db.Task{}, fmt.Errorf("[SERVICE] error creating task: %w", err)
	}

	return task, nil
}

func (s *TaskService) GetTaskInfo(ctx context.Context, id pgtype.UUID) (db.Task, error) {
	log.WithContext(ctx).Infof("Get task for ID %v: ", id)

	task, err := s.db.ReadTask(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Task{}, fmt.Errorf("[SERVICE] task not found: %w", err)
		}
		return db.Task{}, fmt.Errorf("[SERVICE] error getting task: %w", err)
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id pgtype.UUID, text string) (db.Task, error) {
	log.WithContext(ctx).Infof("Update task for ID: %v,\n Text %s", id, text)

	if len(text) <= 0 {
		return db.Task{}, fmt.Errorf("[SERVICE] error - text empty")
	}

	task, err := s.db.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          id,
		Description: text,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Task{}, fmt.Errorf("[SERVICE] task not found: %w", err)
		}
		return db.Task{}, fmt.Errorf("[SERVICE] error updating task: %w", err)
	}

	return task, nil
}

func (s *TaskService) RemoveTask(ctx context.Context, id pgtype.UUID) error {
	log.WithContext(ctx).Infof("Remove task for ID: %v", id)

	if _, err := s.db.ReadTask(ctx, id); err != nil {
		return fmt.Errorf("[SERVICE] task not found: %w", err)
	}

	if err := s.db.RemoveTask(ctx, id); err != nil {
		return fmt.Errorf("[SERVICE] error deleting task: %w", err)
	}

	return nil
}
