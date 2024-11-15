package auth_controller

import (
	"fashora-backend/services/auth_service"
	"fashora-backend/services/user_service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input struct {
		PhoneNumber string     `json:"phone_number" binding:"required"`
		Password    string     `json:"password" binding:"required"`
		UserName    *string    `json:"user_name"` // Optional field
		Birthday    *time.Time `json:"birthday"`  // Optional field, format: "YYYY-MM-DD"
		Address     *string    `json:"address"`   // Optional field
		DeviceID    *string    `json:"device_id"` // Optional field
		Gender      *int       `json:"gender"`    // Optional field, should be 0, 1, or 2
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userWithToken, err := auth_service.Register(user_service.UserRegisterInfo(input))

	if err != nil {
		if err.Error() == "phone number already registered" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already registered"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the new token to the client
	c.JSON(http.StatusCreated, gin.H{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	})
}

func Login(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userWithToken, err := auth_service.Login(input.PhoneNumber, input.Password)

	// TODO use constant and error code? for compare
	if err != nil {
		if err.Error() == "phone number not registered" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Phone number not registered"})
			return
		}

		if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the new token to the client
	c.JSON(http.StatusOK, gin.H{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	})
}
