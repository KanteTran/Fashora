package auth_controller

import (
	"errors"
	"fashora-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CheckPhoneNumberExists(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    nil,
		})
		return
	}

	var user models.Users
	err := models.DB.Where("phone = ?", input.PhoneNumber).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, models.Response{
				Success: true,
				Status:  http.StatusOK,
				Message: "Phone number does not exist",
				Data:    gin.H{"exists": false},
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "Database error",
				Data:    nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
			Status:  http.StatusOK,
			Message: "Phone number exists",
			Data:    gin.H{"exists": true},
		})
	}
}
