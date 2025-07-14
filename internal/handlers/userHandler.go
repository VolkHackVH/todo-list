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

type UserHandler struct {
	Service *service.UserService
}

type UserRequest struct {
	Id       pgtype.UUID `json:"id" binding:"required"`
	Username string      `json:"username" binding:"required,min=3,max=15"`
}

func NewUserHandler(db *db.Queries) *UserHandler {
	return &UserHandler{
		Service: service.NewUserService(db),
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.CreateUser(c.Request.Context(), req.Username)
	if err != nil {
		log.WithContext(c.Request.Context()).Errorf("CreateUser failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created_info": map[string]string{
			"id":       user.ID.String(),
			"username": user.Username,
		},
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	id, ok := helper.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	user, err := h.Service.GetUserInfo(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "User not found",
					"target":  id.String(),
				},
			})
		} else {
			log.WithContext(c.Request.Context()).Errorf("GetUserInfo failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req UserRequest

	id, ok := helper.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.UpdateUserInfo(c.Request.Context(), req.Username, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": "User not found",
					"target":  id.String(),
				},
			})
		} else {
			log.WithContext(c.Request.Context()).Errorf("UpdateUserInfo failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RemoveUser(c *gin.Context) {
	id, ok := helper.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.Service.RemoveUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "User not found",
					"target":  id.String(),
				},
			})
		} else {
			log.WithContext(c.Request.Context()).Errorf("RemoveUser failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("user with ID %s was deleted", id),
	})
}
