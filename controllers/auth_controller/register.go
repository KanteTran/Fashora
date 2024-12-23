package auth_controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"fashora-backend/models"
	"fashora-backend/services/auth_service"
	"fashora-backend/utils"
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
	if !utils.ValidatePhoneOTP(c, input.PhoneNumber) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Phone number is not valid by OTP")
		return
	}
	userWithToken, err := auth_service.Register(models.UserInfo(input))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	responseData := map[string]interface{}{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	}
	utils.SendSuccessResponse(c, http.StatusCreated, "User created successfully", responseData)
}
