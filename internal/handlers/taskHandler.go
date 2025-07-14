package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/VolkHackVH/todo-list.git/internal/db"
	"github.com/VolkHackVH/todo-list.git/internal/helper"
	"github.com/VolkHackVH/todo-list.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

type TaskHandler struct {
	Service *service.TaskService
}

// ? Struct request
type CreateTaskRequest struct {
	Text   string      `json:"description" binding:"required,min=3,max=100"`
	UserID pgtype.UUID `json:"user_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Text string `json:"description" binding:"required,min=3,max=100"`
}

func NewTaskHandler(db *db.Queries) *TaskHandler {
	return &TaskHandler{
		Service: service.NewTaskService(db),
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.Service.CreateTask(c.Request.Context(), req.UserID, req.Text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "User not found",
					"target":  req.UserID.String(),
				}})
		} else {
			log.WithContext(c.Request.Context()).Errorf("CreateTask failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"create_info": task})
}

func (h *TaskHandler) GetTaskInfo(c *gin.Context) {
	var id pgtype.UUID

	id, ok := helper.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	task, err := h.Service.GetTaskInfo(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "Task not found",
					"target":  id.String(),
				}})
		} else {
			log.WithContext(c.Request.Context()).Errorf("GetTaskInfo failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var id pgtype.UUID
	var req UpdateTaskRequest

	id, ok := helper.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.Service.UpdateTask(c.Request.Context(), id, req.Text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": "Task not found",
					"target":  id.String(),
				}})
		} else {
			log.WithContext(c.Request.Context()).Errorf("UpdateTask failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) RemoveTask(c *gin.Context) {
	var id pgtype.UUID

	id, ok := helper.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.Service.RemoveTask(c.Request.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": "Task not found",
					"target":  id.String(),
				}})
		} else {
			log.WithContext(c.Request.Context()).Errorf("RemoveTask failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("task with ID %s was deleted", id)})
}
