package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUIDParam(c *gin.Context, paramName string) (pgtype.UUID, bool) {
	var id pgtype.UUID
	if err := id.Scan(c.Param(paramName)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return pgtype.UUID{}, false
	}
	return id, true
}
