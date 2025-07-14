package handlers

import (
	"github.com/VolkHackVH/todo-list.git/internal/db"
)

type Handlers struct {
	User *UserHandler
	Task *TaskHandler
}

func NewHandler(db *db.Queries) *Handlers {
	return &Handlers{
		User: NewUserHandler(db),
		Task: NewTaskHandler(db),
	}
}
