package auth_controller

import (
	"errors"
	"fashora-backend/models"
	"fashora-backend/services/user_service"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// Update handles updating user information
func Update(c *gin.Context) {
	// Bind the input data
	var input models.UserInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Get the authenticated user
	user, err := getAuthenticatedUser(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Validate token and input data
	if user.Phone != input.PhoneNumber {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid Token")
		return
	}

	// Update the user information in the database
	if err := user_service.UpdateUserByPhoneNumber(input); err != nil {
		handleUpdateError(c, err)
		return
	}

	// Respond with success
	utils.SendSuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

// getAuthenticatedUser extracts the authenticated user from the context
func getAuthenticatedUser(c *gin.Context) (models.Users, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return models.Users{}, errors.New("User not authenticated")
	}

	user, ok := userInterface.(models.Users)
	if !ok {
		return models.Users{}, errors.New("Invalid user type")
	}

	return user, nil
}

// handleUpdateError handles errors during user update
func handleUpdateError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.SendErrorResponse(c, http.StatusNotFound, "User does not exist")
	} else {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
	}
}
