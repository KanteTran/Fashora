package utils

import (
	"fashora-backend/models"

	"github.com/gin-gonic/gin"
)

func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.Response{
		Success: false,
		Status:  statusCode,
		Message: message,
		Data:    nil,
	})
}

func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.Response{
		Success: true,
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}
