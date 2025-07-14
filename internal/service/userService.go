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

type UserService struct {
	db *db.Queries
}

func NewUserService(query *db.Queries) *UserService {
	return &UserService{db: query}
}

func (s *UserService) CreateUser(ctx context.Context, username string) (db.User, error) {
	log.WithFields(log.Fields{
		"username": username,
	}).Info("Created user")

	if len(username) <= 0 {
		return db.User{}, fmt.Errorf("[SERVICE] username empty")
	}

	user, err := s.db.CreateUser(ctx, username)
	if err != nil {
		return db.User{}, fmt.Errorf("[SERVICE] error creation user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, id pgtype.UUID) (db.User, error) {
	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting user")

	user, err := s.db.ReadUserInfo(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, fmt.Errorf("[SERVICE] ID not found: %w", pgx.ErrNoRows)
		}
		return db.User{}, fmt.Errorf("[SERVICE] error getting user info: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, username string, id pgtype.UUID) (db.User, error) {
	log.WithFields(log.Fields{
		"id":       id,
		"username": username,
	}).Info("Updating user info")

	user, err := s.db.UpdateUser(ctx, db.UpdateUserParams{
		ID:       id,
		Username: username,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, fmt.Errorf("[SERVICE] id not found: %w", pgx.ErrNoRows)
		}
		return db.User{}, fmt.Errorf("[SERVICE] error updating user info: %w", err)
	}

	return user, nil
}

func (s *UserService) RemoveUser(ctx context.Context, id pgtype.UUID) error {
	log.WithFields(log.Fields{
		"id": id,
	}).Info("Deleting user")

	if _, err := s.db.ReadUserInfo(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: user not found", pgx.ErrNoRows)
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if err := s.db.RemoveUser(ctx, id); err != nil {
		return fmt.Errorf("[SERVICE] error remove user: %w", err)
	}

	return nil
}
