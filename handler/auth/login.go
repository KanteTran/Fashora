package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fashora-backend/services/auth_service"
	"fashora-backend/utils"
)

type LoginInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// Login handles user login
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	userWithToken, err := auth_service.Login(input.PhoneNumber, input.Password)
	if err != nil {
		handleLoginError(c, err)
		return
	}

	responseData := map[string]interface{}{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", responseData)
}

// handleLoginError handles login-specific errors
func handleLoginError(c *gin.Context, err error) {
	switch err.Error() {
	case "phone number not registered":
		utils.SendErrorResponse(c, http.StatusNotFound, "Phone number not registered")
	case "invalid password":
		utils.SendErrorResponse(c, http.StatusNotFound, "Invalid password")
	default:
		utils.SendErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
