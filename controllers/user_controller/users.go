package user_controller

import (
	"fashora-backend/services/user_service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateUser(c *gin.Context) {
	var input struct {
		PhoneNumber string     `json:"phone_number" binding:"required"`
		Password    *string    `json:"password"`  // Optional, hashed if provided
		UserName    *string    `json:"user_name"` // Optional
		Birthday    *time.Time `json:"birthday"`  // Optional, format: "YYYY-MM-DD"
		Address     *string    `json:"address"`   // Optional
		DeviceID    *string    `json:"device_id"` // Optional
		Gender      *int       `json:"gender"`    // Optional, should be 0, 1, or 2
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := user_service.UpdateUserByPhoneNumber(user_service.UserUpdateInfo(input)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not exist"})
		}
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
