package auth_controller

import (
	"fashora-backend/services/auth_service"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// Login handles user login
func Login(c *gin.Context) {
	// Bind input to LoginInput struct
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Authenticate the user
	userWithToken, err := auth_service.Login(input.PhoneNumber, input.Password)
	if err != nil {
		handleLoginError(c, err)
		return
	}

	// Return success response with user data
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
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Phone number not registered")
	case "invalid password":
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid password")
	default:
		utils.SendErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred")
	}
}
