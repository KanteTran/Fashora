package auth_controller

import (
	"fashora-backend/models"
	"fashora-backend/services/auth_service"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RegisterInput struct {
	PhoneNumber string     `json:"phone_number" binding:"required"`
	Password    string     `json:"password" binding:"required"`
	UserName    *string    `json:"user_name"`
	Birthday    *time.Time `json:"birthday"`
	Address     *string    `json:"address"`
	DeviceID    *string    `json:"device_id"`
	Gender      *int       `json:"gender"` // 0: male, 1: female, 2: other
}

// Register handles user registration
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if !utils.ValidatePhoneNumber(input.PhoneNumber) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Phone number is invalid")
		return
	}

	userWithToken, err := auth_service.Register(models.UserInfo(input))
	if err != nil {
		handleRegistrationError(c, err)
		return
	}

	responseData := map[string]interface{}{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	}
	utils.SendSuccessResponse(c, http.StatusCreated, "User created successfully", responseData)
}

// handleRegistrationError handles errors returned from the registration service
func handleRegistrationError(c *gin.Context, err error) {
	switch err.Error() {
	case "phone number already registered":
		utils.SendErrorResponse(c, http.StatusBadRequest, "Phone number already registered")
	default:
		utils.SendErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
